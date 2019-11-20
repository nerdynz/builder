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

// FormField Struct
type FormField struct {
	FormFieldID  int              `db:"form_field_id" json:"FormFieldID"`
	Name         string           `db:"name" json:"Name"`
	FieldName    string           `db:"field_name" json:"FieldName"`
	Type         string           `db:"type" json:"Type"`
	IsRequired   bool             `db:"is_required" json:"IsRequired"`
	FormID       int              `db:"form_id" json:"FormID"`
	DateCreated  time.Time        `db:"date_created" json:"DateCreated"`
	DateModified time.Time        `db:"date_modified" json:"DateModified"`
	FieldOptions FormFieldOptions `json:"FieldOptions"`
	UUID         string           `db:"uuid" json:"UUID"`
}

var formFieldHelperGlobal *formFieldHelper

type FormFields []*FormField

type formFieldHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func FormFieldHelper() *formFieldHelper {
	if formFieldHelperGlobal == nil {
		formFieldHelperGlobal = newFormFieldHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return formFieldHelperGlobal
}

func newFormFieldHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *formFieldHelper {
	helper := &formFieldHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"form_field_id", "name", "field_name", "type", "is_required", "form_id", "date_created", "date_modified", "uuid"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "date_created, date_modified"

	return helper
}

func (h *formFieldHelper) New() *FormField {
	record := &FormField{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *formFieldHelper) NewFromRequest(req *http.Request) (*FormField, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *formFieldHelper) LoadAndUpdateFromRequest(req *http.Request) (*FormField, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.FormFieldID <= 0 {
		return nil, errors.New("The  failed to load because FormFieldID was not found in the request.")
	}

	d, err := time.ParseDuration("30s")
	if err != nil {
		return nil, err
	}
	slightlyAheadDateModified := newRecord.DateModified.Add(d)
	record, err := h.Load(newRecord.FormFieldID)
	if record.DateModified.After(slightlyAheadDateModified) {
		errMsg := "This FormField record has been modified recently. Please refresh the browser to load the latest changes."
		errMsg += "DEVDETAILS: The FormField record failed to update because the DateModified value in the database is more recent then DateModified value on the request.\n"
		errMsg += "request: [" + newRecord.DateModified.String() + "]\n"
		errMsg += "database: [" + record.DateModified.String() + "]\n"
		return nil, errors.New(errMsg)
	}

	newRecord.FormFieldID = record.FormFieldID // this shouldn't have changed
	newRecord.DateCreated = record.DateCreated // nor should this.

	return newRecord, nil
}

func (h *formFieldHelper) UpdateFromRequest(req *http.Request, record *FormField) error {
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

func (h *formFieldHelper) All() (FormFields, error) {
	var records FormFields
	err := h.DB.Select("*").
		From("form_field").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (FormFields, error) {
	var records FormFields
	err := h.DB.Select("*").
		From("form_field").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	for _, record := range records {
		options, err := FormFieldOptionHelper().Where("form_field_id = $1", record.FormFieldID)
		if err != nil && err.Error() != NoRows {
			return nil, err
		}
		record.FieldOptions = options
	}

	return records, nil
}

func (h *formFieldHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*FormField, error) {
	var record FormField

	err := h.DB.Select("*").
		From("form_field").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *formFieldHelper) Paged(pageNum int, itemsPerPage int) (FormFields, error) {
	var records FormFields
	records, err := h.PagedBy(pageNum, itemsPerPage, "date_created") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string) (FormFields, error) {
	i := sort.SearchStrings(h.fieldNames, orderByFieldName)
	// check the orderby exists within the fields as this could be an easy sql injection hole.
	if !(i < len(h.fieldNames) && h.fieldNames[i] == orderByFieldName) { // NOT
		return nil, errors.New("field name [" + orderByFieldName + "]  isn't a valid field name")
	}

	var records FormFields
	err := h.DB.Select("*").
		From("form_field").
		OrderBy(orderByFieldName).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldHelper) Load(id int) (*FormField, error) {
	record := &FormField{}
	err := h.DB.
		Select("*").
		From("form_field").
		Where("form_field_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *formFieldHelper) Save(record *FormField) error {
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

func (h *formFieldHelper) SaveMany(records FormFields) error {
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

func (h *formFieldHelper) save(record *FormField) error {
	err := h.DB.
		Upsert("form_field").
		Columns("name", "field_name", "type", "is_required", "form_id", "date_created", "date_modified", "uuid").
		Values(record.Name, record.FieldName, record.Type, record.IsRequired, record.FormID, record.DateCreated, record.DateModified, record.UUID).
		Where("form_field_id=$1", record.FormFieldID).
		Returning("form_field_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	//if len(record.DayItems) > 0 {
	//	for _, dayItem := range record.DayItems {
	//		dayItem.DayID = record.DayID // may have just been set
	//	}
	//	dayItemHelper := NewDayItemHelper(h.DB, h.Cache)
	//	err := dayItemHelper.SaveMany(record.DayItems)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

// Validate a record
func (h *formFieldHelper) Validate(record *FormField) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *formFieldHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("form_field").
		Where("form_field_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
