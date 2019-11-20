package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/avct/uasurfer"
	"github.com/gorilla/schema"

	"strings"

	"sort"

	runner "github.com/nerdynz/dat/sqlx-runner"
	validator "gopkg.in/go-playground/validator.v9"
	redis "github.com/go-redis/redis"
)

// Analytics Struct
type Analytics struct {
	AnalyticsID  int       `db:"analytics_id" json:"AnalyticsID"`
	Page         string    `db:"page" json:"Page"`
	ScreenWidth  int       `db:"screen_width" json:"ScreenWidth"`
	ScreenHeight int       `db:"screen_height" json:"ScreenHeight"`
	UserAgent    string    `db:"user_agent" json:"UserAgent"`
	VpWidth      int       `db:"vp_width" json:"VpWidth"`
	VpHeight     int       `db:"vp_height" json:"VpHeight"`
	UniqueID     float64   `db:"unique_id" json:"UniqueID"`
	DateCreated  time.Time `db:"date_created" json:"DateCreated"`
	DateModified time.Time `db:"date_modified" json:"DateModified"`
	Browser      string    `db:"browser" json:"Browser"`
	Device       string    `db:"device" json:"Device"`
	Version      int       `db:"version" json:"Version"`
}

func (a *Analytics) SetUniqueID(str string) {
	if fl, err := strconv.ParseFloat(str, 64); err == nil {
		a.UniqueID = fl
	}
}

func (a *Analytics) SetUaInfo() {
	if a.UserAgent != "" {
		det := uasurfer.Parse(a.UserAgent)
		a.Browser = strings.TrimPrefix(det.Browser.Name.String(), "Browser")
		a.Device = strings.TrimPrefix(det.DeviceType.String(), "Device")
		a.Version = det.Browser.Version.Major
	}
}

var analyticsHelperGlobal *analyticsHelper

type Analyticss []*Analytics

type analyticsHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func AnalyticsHelper() *analyticsHelper {
	if analyticsHelperGlobal == nil {
		analyticsHelperGlobal = newAnalyticsHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return analyticsHelperGlobal
}

func newAnalyticsHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *analyticsHelper {
	helper := &analyticsHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"analytics_id", "page", "screen_width", "screen_height", "user_agent", "vp_width", "vp_height", "unique_id", "date_created", "date_modified"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "date_created, date_modified"

	return helper
}

func (h *analyticsHelper) New() *Analytics {
	record := &Analytics{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *analyticsHelper) NewFromRequest(req *http.Request) (*Analytics, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *analyticsHelper) LoadAndUpdateFromRequest(req *http.Request) (*Analytics, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.AnalyticsID <= 0 {
		return nil, errors.New("The  failed to load because AnalyticsID was not found in the request.")
	}

	return newRecord, nil
}

func (h *analyticsHelper) UpdateFromRequest(req *http.Request, record *Analytics) error {
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

func (h *analyticsHelper) All() (Analyticss, error) {
	var records Analyticss
	err := h.DB.Select(strings.Join(h.fieldNames, ",")).
		From("analytics").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *analyticsHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (Analyticss, error) {
	var records Analyticss
	err := h.DB.Select(strings.Join(h.fieldNames, ",")).
		From("analytics").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}
func (h *analyticsHelper) SQL(sql string, args ...interface{}) (Analyticss, error) {
	var records Analyticss
	err := h.DB.SQL(sql, args...).
		QueryStructs(&records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (h *analyticsHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*Analytics, error) {
	var record Analytics

	err := h.DB.Select(strings.Join(h.fieldNames, ",")).
		From("analytics").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *analyticsHelper) Load(id int) (*Analytics, error) {
	record := &Analytics{}
	err := h.DB.
		Select(strings.Join(h.fieldNames, ",")).
		From("analytics").
		Where("analytics_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *analyticsHelper) Save(record *Analytics) error {
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

func (h *analyticsHelper) SaveMany(records Analyticss) error {
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

func (h *analyticsHelper) save(record *Analytics) error {
	if record.UniqueID == 0 {
		return errors.New("invalid unique can't save")
	}
	err := h.DB.
		Upsert("analytics").
		Columns("page", "screen_width", "screen_height", "user_agent", "vp_width", "vp_height", "unique_id", "date_created", "date_modified", "device", "version", "browser").
		Values(record.Page, record.ScreenWidth, record.ScreenHeight, record.UserAgent, record.VpWidth, record.VpHeight, record.UniqueID, record.DateCreated, record.DateModified, record.Device, record.Version, record.Browser).
		Where("analytics_id=$1", record.AnalyticsID).
		Returning("analytics_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

// Validate a record
func (h *analyticsHelper) Validate(record *Analytics) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *analyticsHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("analytics").
		Where("analytics_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
