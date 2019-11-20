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

// FormFieldOption Struct
type FormFieldOption struct {
	FormFieldOptionID int       `db:"form_field_option_id" json:"FormFieldOptionID"`
	Label             string    `db:"label" json:"Label"`
	Value             string    `db:"value" json:"Value"`
	DateCreated       time.Time `db:"date_created" json:"DateCreated"`
	DateModified      time.Time `db:"date_modified" json:"DateModified"`
	UUID              string    `db:"uuid" json:"UUID"`
	FormFieldID       int       `db:"form_field_id" json:"FormFieldID"`
	IsDefault         bool      `db:"is_default" json:"IsDefault"`
}

var formFieldOptionHelperGlobal *formFieldOptionHelper

type FormFieldOptions []*FormFieldOption

type formFieldOptionHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func FormFieldOptionHelper() *formFieldOptionHelper {
	if formFieldOptionHelperGlobal == nil {
		formFieldOptionHelperGlobal = newFormFieldOptionHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return formFieldOptionHelperGlobal
}

func newFormFieldOptionHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *formFieldOptionHelper {
	helper := &formFieldOptionHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"form_field_option_id", "label", "value", "date_created", "date_modified", "uuid", "form_field_id", "is_default"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "date_created, date_modified"

	return helper
}

func (h *formFieldOptionHelper) New() *FormFieldOption {
	record := &FormFieldOption{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *formFieldOptionHelper) NewFromRequest(req *http.Request) (*FormFieldOption, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *formFieldOptionHelper) LoadAndUpdateFromRequest(req *http.Request) (*FormFieldOption, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.FormFieldOptionID <= 0 {
		return nil, errors.New("The  failed to load because FormFieldOptionID was not found in the request.")
	}

	d, err := time.ParseDuration("30s")
	if err != nil {
		return nil, err
	}
	slightlyAheadDateModified := newRecord.DateModified.Add(d)
	record, err := h.Load(newRecord.FormFieldOptionID)
	if record.DateModified.After(slightlyAheadDateModified) {
		errMsg := "This FormFieldOption record has been modified recently. Please refresh the browser to load the latest changes."
		errMsg += "DEVDETAILS: The FormFieldOption record failed to update because the DateModified value in the database is more recent then DateModified value on the request.\n"
		errMsg += "request: [" + newRecord.DateModified.String() + "]\n"
		errMsg += "database: [" + record.DateModified.String() + "]\n"
		return nil, errors.New(errMsg)
	}

	newRecord.FormFieldOptionID = record.FormFieldOptionID // this shouldn't have changed
	newRecord.DateCreated = record.DateCreated             // nor should this.

	return newRecord, nil
}

func (h *formFieldOptionHelper) UpdateFromRequest(req *http.Request, record *FormFieldOption) error {
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

func (h *formFieldOptionHelper) All() (FormFieldOptions, error) {
	var records FormFieldOptions
	err := h.DB.Select("*").
		From("form_field_option").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldOptionHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (FormFieldOptions, error) {
	var records FormFieldOptions
	err := h.DB.Select("*").
		From("form_field_option").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldOptionHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*FormFieldOption, error) {
	var record FormFieldOption

	err := h.DB.Select("*").
		From("form_field_option").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *formFieldOptionHelper) Paged(pageNum int, itemsPerPage int) (FormFieldOptions, error) {
	var records FormFieldOptions
	records, err := h.PagedBy(pageNum, itemsPerPage, "date_created") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldOptionHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string) (FormFieldOptions, error) {
	i := sort.SearchStrings(h.fieldNames, orderByFieldName)
	// check the orderby exists within the fields as this could be an easy sql injection hole.
	if !(i < len(h.fieldNames) && h.fieldNames[i] == orderByFieldName) { // NOT
		return nil, errors.New("field name [" + orderByFieldName + "]  isn't a valid field name")
	}

	var records FormFieldOptions
	err := h.DB.Select("*").
		From("form_field_option").
		OrderBy(orderByFieldName).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *formFieldOptionHelper) Load(id int) (*FormFieldOption, error) {
	record := &FormFieldOption{}
	err := h.DB.
		Select("*").
		From("form_field_option").
		Where("form_field_option_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *formFieldOptionHelper) Save(record *FormFieldOption) error {
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

func (h *formFieldOptionHelper) SaveMany(records FormFieldOptions) error {
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

func (h *formFieldOptionHelper) save(record *FormFieldOption) error {
	err := h.DB.
		Upsert("form_field_option").
		Columns("label", "value", "date_created", "date_modified", "uuid", "form_field_id", "is_default").
		Values(record.Label, record.Value, record.DateCreated, record.DateModified, record.UUID, record.FormFieldID, record.IsDefault).
		Where("form_field_option_id=$1", record.FormFieldOptionID).
		Returning("form_field_option_id").
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
func (h *formFieldOptionHelper) Validate(record *FormFieldOption) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *formFieldOptionHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("form_field_option").
		Where("form_field_option_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
