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
	"github.com/nerdynz/security"
	validator "gopkg.in/go-playground/validator.v9"
	redis "github.com/go-redis/redis"
)

var personHelperGlobal *personHelper

// Person Record
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
	ULID         string    `db:"ulid" json:"ULID"`
}

type People []*Person

func (h *personHelper) beforeSave(record *Person) (err error) {
	if record.DateCreated.IsZero() {
		record.DateCreated = time.Now()
	}
	record.DateModified = time.Now()
	if record.ULID == "" {
		record.ULID = security.ULID()
	}

	validationErr := h.validate(record)
	if validationErr != nil {
		return validationErr
	}
	return err
}

func (h *personHelper) afterSave(record *Person) (err error) {
	return err
}

// GENERATED CODE - Leave the below code alone
type personHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	structDecoder *schema.Decoder
	validator     *validator.Validate
	fieldNames    []string
	orderBy       string
}

func PersonHelper() *personHelper {
	if personHelperGlobal == nil {
		personHelperGlobal = newPersonHelper(modelDB, modelCache, modelDecoder, modelValidator)
	}
	return personHelperGlobal
}

func newPersonHelper(db *runner.DB, redis *redis.Client, d *schema.Decoder, v *validator.Validate) *personHelper {
	helper := &personHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.structDecoder = d
	helper.validator = v

	// Fields
	fieldnames := []string{"person_id", "name", "email", "password", "phone", "role", "picture", "date_created", "date_modified", "ulid"}
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

func (h *personHelper) FromRequest(req *http.Request) (*Person, error) {
	record := h.New()
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// working with json
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(record)
		if err != nil {
			return nil, err
		}
	} else {
		// working with form values
		err := req.ParseForm()
		if err != nil {
			return nil, err
		}

		err = h.structDecoder.Decode(record, req.PostForm)
		if err != nil {
			return nil, err
		}
	}
	return record, nil
}

func (h *personHelper) Load(id int) (*Person, error) {
	record, err := h.One("person_id = $1", id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *personHelper) All() (People, error) {
	var records People
	err := h.DB.Select("*").
		From("person").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *personHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (People, error) {
	var records People
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

	var records People
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

func (h *personHelper) Save(record *Person) error {
	return h.save(record)
}

func (h *personHelper) SaveMany(records People) error {
	for _, record := range records {
		// everything is validated so now re loop and do the actual saving... this should probably be a tx that can just rollback
		err := h.save(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *personHelper) save(record *Person) (err error) {
	err = h.beforeSave(record)
	if err != nil {
		return err
	}
	cols := []string{"name", "email", "password", "phone", "role", "picture", "date_created", "date_modified", "ulid"}
	vals := []interface{}{record.Name, record.Email, record.Password, record.Phone, record.Role, record.Picture, record.DateCreated, record.DateModified, record.ULID}
	if record.PersonID > 0 {
		// UPDATE
		b := h.DB.Update("person")
		for i := range cols {
			b.Set(cols[i], vals[i])
		}
		b.Where("person_id = $1", record.PersonID)
		b.Returning("person_id")
		err = b.QueryStruct(record)
	} else {
		// INSERT
		err = h.DB.
			InsertInto("person").
			Columns(cols...).
			Values(vals...).
			Returning("person_id").
			QueryStruct(record)
	}
	if err != nil {
		return err
	}
	err = h.afterSave(record)
	return err
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

func (h *personHelper) validate(record *Person) (err error) {
	validationErrors := h.validator.Struct(record)
	if validationErrors != nil {
		errMessage := ""
		for _, err := range err.(validator.ValidationErrors) {
			errMessage += err.Kind().String() + " validation Error on field " + err.Field()
		}
		if errMessage != "" {
			err = errors.New(errMessage)
		}
	}
	return err
}
