package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/schema"

	"strings"

	"sort"

	runner "github.com/bkono/dat/sqlx-runner"
	validator "gopkg.in/go-playground/validator.v9"
	redis "gopkg.in/redis.v5"
)

// Setting Struct
type Setting struct {
	SettingID      int       `db:"setting_id" json:"SettingID"`
	FooterTextHTML string    `db:"footer_text_html" json:"FooterTextHTML"`
	DateCreated    time.Time `db:"date_added" json:"DateAdded"`
	DateModified   time.Time `db:"date_modified" json:"DateModified"`
}

var settingHelperGlobal *settingHelper

type Settings []*Setting

type settingHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func SettingHelper() *settingHelper {
	if settingHelperGlobal == nil {
		settingHelperGlobal = newSettingHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return settingHelperGlobal
}

func newSettingHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *settingHelper {
	helper := &settingHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"setting_id", "footer_text_html", "date_modified", "date_added"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames

	helper.orderBy = "date_created, date_modified"
	return helper
}

func (h *settingHelper) New() *Setting {
	record := &Setting{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *settingHelper) NewFromRequest(req *http.Request) (*Setting, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *settingHelper) LoadAndUpdateFromRequest(req *http.Request) (*Setting, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.SettingID <= 0 {
		return nil, errors.New("The  failed to load because SettingID was not found in the request.")
	}
	return newRecord, nil
}

func (h *settingHelper) UpdateFromRequest(req *http.Request, record *Setting) error {
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

func (h *settingHelper) All() (Settings, error) {
	var records Settings
	err := h.DB.Select("*").
		From("setting").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *settingHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (Settings, error) {
	var records Settings
	err := h.DB.Select("*").
		From("setting").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *settingHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*Setting, error) {
	var record Setting

	err := h.DB.Select("*").
		From("setting").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *settingHelper) Paged(pageNum int, itemsPerPage int) (*PagedData, error) {
	pd, err := h.PagedBy(pageNum, itemsPerPage, "date_created", "") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (h *settingHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string, direction string) (*PagedData, error) {
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

	var records Settings
	err := h.DB.Select("*").
		From("setting").
		OrderBy(orderByFieldName + " " + direction).
		Offset(uint64((pageNum - 1) * itemsPerPage)).
		Limit(uint64(itemsPerPage)).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	count := 0
	h.DB.SQL(`select count(setting_id) from setting`).QueryStruct(&count)
	return NewPagedData(records, orderByFieldName, direction, itemsPerPage, pageNum, count), nil
}

func (h *settingHelper) Load(id int) (*Setting, error) {
	record := &Setting{}
	err := h.DB.
		Select("*").
		From("setting").
		Where("setting_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *settingHelper) Save(record *Setting) error {
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

func (h *settingHelper) SaveMany(records Settings) error {
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

func (h *settingHelper) save(record *Setting) error {
	err := h.DB.
		Upsert("setting").
		Columns("footer_text_html", "date_modified", "date_added").
		Values(record.FooterTextHTML, record.DateModified, record.DateCreated).
		Where("setting_id=$1", record.SettingID).
		Returning("setting_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

// Validate a record
func (h *settingHelper) Validate(record *Setting) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *settingHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("setting").
		Where("setting_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
