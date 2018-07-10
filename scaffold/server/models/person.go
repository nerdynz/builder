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
	redis "gopkg.in/redis.v5"
)

// Person Struct
type Person struct {
	PersonID     int       `db:"person_id" json:"PersonID"`
	Name         string    `db:"name" json:"Name"`
	Email        string    `db:"email" json:"Email"`
	Password     string    `db:"password" json:"Password"`
	Phone        string    `db:"phone" json:"Phone"`
	Role         string    `db:"role" json:"Role"`
	Picture      string    `db:"picture" json:"Picture"`
	DateCreated  time.Time `db:"date_created" json:"DateCreated"`
	DateModified time.Time `db:"date_modified" json:"DateModified"`
}

var personHelperGlobal *personHelper

type Persons []*Person

type personHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func PersonHelper() *personHelper {
	if personHelperGlobal == nil {
		personHelperGlobal = newPersonHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return personHelperGlobal
}

func newPersonHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *personHelper {
	helper := &personHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"person_id", "name", "email", "password", "phone", "role", "picture", "date_created", "date_modified"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "date_created, date_modified"

	return helper
}

func (h *personHelper) New() *Person {
	record := &Person{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *personHelper) NewFromRequest(req *http.Request) (*Person, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *personHelper) LoadAndUpdateFromRequest(req *http.Request) (*Person, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.PersonID <= 0 {
		return nil, errors.New("The  failed to load because PersonID was not found in the request.")
	}

	return newRecord, nil
}

func (h *personHelper) UpdateFromRequest(req *http.Request, record *Person) error {
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

func (h *personHelper) All() (Persons, error) {
	var records Persons
	err := h.DB.Select("*").
		From("person").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *personHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (Persons, error) {
	var records Persons
	err := h.DB.Select("*").
		From("person").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *personHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*Person, error) {
	var record Person

	err := h.DB.Select("*").
		From("person").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *personHelper) Paged(pageNum int, itemsPerPage int) (*PagedData, error) {
	pd, err := h.PagedBy(pageNum, itemsPerPage, "date_created", "") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (h *personHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string, direction string) (*PagedData, error) {
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

	var records Persons
	err := h.DB.Select("*").
		From("person").
		OrderBy(orderByFieldName + " " + direction).
		Offset(uint64((pageNum - 1) * itemsPerPage)).
		Limit(uint64(itemsPerPage)).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	count := 0
	h.DB.SQL(`select count(person_id) from person`).QueryStruct(&count)
	return NewPagedData(records, orderByFieldName, direction, itemsPerPage, pageNum, count), nil
}

func (h *personHelper) Load(id int) (*Person, error) {
	record := &Person{}
	err := h.DB.
		Select("*").
		From("person").
		Where("person_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *personHelper) Save(record *Person) error {
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

func (h *personHelper) SaveMany(records Persons) error {
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

func (h *personHelper) save(record *Person) error {
	err := h.DB.
		Upsert("person").
		Columns("name", "email", "password", "phone", "role", "picture", "date_created", "date_modified").
		Values(record.Name, record.Email, record.Password, record.Phone, record.Role, record.Picture, record.DateCreated, record.DateModified).
		Where("person_id=$1", record.PersonID).
		Returning("person_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

// Validate a record
func (h *personHelper) Validate(record *Person) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *personHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("person").
		Where("person_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
