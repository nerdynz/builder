package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/schema"

	"strings"

	"sort"

	runner "github.com/nerdynz/dat/sqlx-runner"
	validator "gopkg.in/go-playground/validator.v9"
	redis "github.com/go-redis/redis"
)

// Form Struct
type Form struct {
	FormID        int        `db:"form_id" json:"FormID"`
	Name          string     `db:"name" json:"Name"`
	DateCreated   time.Time  `db:"date_created" json:"DateCreated"`
	DateModified  time.Time  `db:"date_modified" json:"DateModified"`
	UUID          string     `db:"uuid" json:"UUID"`
	SubmitText    string     `db:"submit_text" json:"SubmitText"`
	ThanksMessage string     `db:"thanks_message" json:"ThanksMessage"`
	FormURL       string     `db:"form_url" json:"FormURL"`
	FormFields    FormFields `json:"FormFields"`
	SpecialForm   string     `db:"special_form" json:"SpecialForm"`
}

var formHelperGlobal *formHelper

type Forms []*Form

type formHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func FormHelper() *formHelper {
	if formHelperGlobal == nil {
		formHelperGlobal = newFormHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return formHelperGlobal
}

func newFormHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *formHelper {
	helper := &formHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"form_id", "name", "date_created", "date_modified", "uuid", "submit_text", "thanks_message", "form_url", "special_form"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "date_created, date_modified"

	return helper
}

func (h *formHelper) New() *Form {
	record := &Form{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *formHelper) NewFromRequest(req *http.Request) (*Form, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *formHelper) LoadAndUpdateFromRequest(req *http.Request) (*Form, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.FormID <= 0 {
		return nil, errors.New("The  failed to load because FormID was not found in the request.")
	}

	return newRecord, nil
}

func (h *formHelper) UpdateFromRequest(req *http.Request, record *Form) error {
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// working with json
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(record)
		if err != nil {
			return err
		}
	} else {
		// working with form values
		err := req.ParseForm()
		if err != nil {
			return err
		}

		err = h.structDecoder.Decode(record, req.PostForm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *formHelper) All() (Forms, error) {
	var records Forms
	err := h.DB.Select("*").
		From("form").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (Forms, error) {
	var records Forms
	err := h.DB.Select("*").
		From("form").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*Form, error) {
	var record Form

	err := h.DB.Select("*").
		From("form").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *formHelper) Paged(pageNum int, itemsPerPage int) (*PagedData, error) {
	pd, err := h.PagedBy(pageNum, itemsPerPage, "date_created", "") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (h *formHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string, direction string) (*PagedData, error) {
	if orderByFieldName == "" || orderByFieldName == "default" {
		// we only want the first field name
		orderByFieldName = strings.Split(h.orderBy, ",")[0]
		orderByFieldName = strings.Trim(orderByFieldName, " ")
	}
	i := sort.SearchStrings(h.fieldNames, orderByFieldName)
	// check the orderby exists within the fields as this could be an easy sql injection hole.
	if !(i < len(h.fieldNames) && h.fieldNames[i] == orderByFieldName) { // NOT
		return nil, errors.New("field name [" + orderByFieldName + "]  isn't a valid field name")
	}

	if !(direction == "asc" || direction == "desc" || direction == "") {
		return nil, errors.New("direction isn't valid")
	}

	var records Forms
	err := h.DB.Select("*").
		From("form").
		OrderBy(orderByFieldName + " " + direction).
		Offset(uint64((pageNum - 1) * itemsPerPage)).
		Limit(uint64(itemsPerPage)).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	count := 0
	h.DB.SQL(`select count(form_id) from form`).QueryStruct(&count)
	return NewPagedData(records, orderByFieldName, direction, itemsPerPage, pageNum, count), nil
}

func (h *formHelper) Load(id int) (*Form, error) {
	record := &Form{}
	err := h.DB.
		Select("*").
		From("form").
		Where("form_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *formHelper) LoadFullForm(formName string) (*Form, error) {
	record, err := h.One("name = $1", formName)
	if err != nil {
		return nil, err
	}
	return h.loadFullForm(record)
}

func (h *formHelper) LoadFullFormByID(formID int) (*Form, error) {
	record, err := h.One("form_id = $1", formID)
	if err != nil {
		return nil, err
	}
	return h.loadFullForm(record)
}

func (h *formHelper) loadFullForm(record *Form) (*Form, error) {
	formFields, err := FormFieldHelper().Where("form_id = $1", record.FormID)
	if err != nil {
		return nil, err
	}
	record.FormFields = formFields
	h.Cache.Set("form-"+record.Name, record, 24*time.Hour)
	return record, nil
}

func (h *formHelper) Save(record *Form) error {
	// date created always should be set, if its null just set it to now.
	if record.DateCreated.IsZero() {
		record.DateCreated = time.Now()
	}

	// was just modified
	record.DateModified = time.Now()

	// check validation
	_, err := h.Validate(record)
	if err != nil {
		return err
	}

	err = h.save(record)
	if err != nil {
		return err
	}

	return err
}

func (h *formHelper) SaveMany(records Forms) error {
	for _, record := range records {
		// date created always should be set, if its null just set it to now.
		if record.DateCreated.IsZero() {
			record.DateCreated = time.Now()
		}

		// was just modified
		record.DateModified = time.Now()

		// check validation
		_, err := h.Validate(record)
		if err != nil {
			return err
		}
	}

	for _, record := range records {
		// everything is validated so now re loop and do the actual saving... this should probably be a tx that can just rollback
		err := h.save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *formHelper) save(record *Form) error {
	err := h.DB.
		Upsert("form").
		Columns("name", "date_created", "date_modified", "uuid", "submit_text", "thanks_message", "form_url", "special_form").
		Values(record.Name, record.DateCreated, record.DateModified, record.UUID, record.SubmitText, record.ThanksMessage, record.FormURL, record.SpecialForm).
		Where("form_id=$1", record.FormID).
		Returning("form_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

// Validate a record
func (h *formHelper) Validate(record *Form) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *formHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("form").
		Where("form_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
