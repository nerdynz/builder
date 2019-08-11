package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/schema"

	"strings"

	"sort"

	runner "github.com/nerdynz/dat/sqlx-runner"
	redis "gopkg.in/redis.v5"
)

// ImageMeta Struct
type ImageMeta struct {
	ImageMetaID     int       `db:"image_meta_id" json:"imageMetaID"`
	Name            string    `db:"name" json:"name"`
	Original        string    `db:"original" json:"originalName"`
	CSSHeight       float64   `db:"css_height" json:"cssHeight"`
	CSSWidth        float64   `db:"css_width" json:"cssWidth"`
	CSSLeft         float64   `db:"css_left" json:"cssLeft"`
	CSSTop          float64   `db:"css_top" json:"cssTop"`
	ContainerHeight float64   `db:"container_height" json:"containerHeight"`
	ContainerWidth  float64   `db:"container_width" json:"containerWidth"`
	Width           int       `json:"width"`
	Height          int       `json:"height"`
	Left            int       `json:"left"`
	Top             int       `json:"top"`
	OriginalHeight  int       `json:"imageOriginalHeight"`
	OriginalWidth   int       `json:"imageOriginalWidth"`
	DateCreated     time.Time `db:"date_created" json:"dateCreated"`
	DateModified    time.Time `db:"date_modified" json:"dateModified"`
	UniqueID        string    `db:"unique_id" json:"uniqueID"`
	Ext             string    `json:"ext"`
	OriginalExt     string    `json:"originalExt"`
	Data            string    `json:"data"`
	OriginalData    string    `json:"original"`
	IsExisting      bool      `json:"isExisting"`
	OldFileName     string    `json:"oldFilename"`
}

// Bytes returns the base64 data passed from request as bytes for the new image
func (meta *ImageMeta) Bytes() ([]byte, error) {
	return getBytesFromBase64(meta.Data, meta.OriginalExt) // always use original ext
}

// OriginalBytes returns the base64 data passed from request as bytes for the new image
func (meta *ImageMeta) OriginalBytes() ([]byte, error) {
	return getBytesFromBase64(meta.OriginalData, meta.OriginalExt) // always use original ext
}

func getBytesFromBase64(data string, ext string) ([]byte, error) {
	spliter := "data:image/jpeg;base64,"
	if ext == "png" {
		spliter = "data:image/png;base64,"
	}
	d := data[len(spliter):]
	return base64.StdEncoding.DecodeString(d)
}

func (meta *ImageMeta) IsConvert() bool {
	return (meta.Ext != meta.OriginalExt)
}

var imageMetaHelperGlobal *imageMetaHelper

type ImageMetas []*ImageMeta

type imageMetaHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func ImageMetaHelper() *imageMetaHelper {
	if imageMetaHelperGlobal == nil {
		imageMetaHelperGlobal = newImageMetaHelper(modelDB, modelCache, modelDecoder)
	}
	return imageMetaHelperGlobal
}

func newImageMetaHelper(db *runner.DB, redis *redis.Client, structDecoder *schema.Decoder) *imageMetaHelper {
	helper := &imageMetaHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"image_meta_id", "name", "original", "css_height", "css_width", "container_height", "container_width", "css_left", "css_top", "date_created", "date_modified", "unique_id"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "date_created, date_modified"

	return helper
}

func (h *imageMetaHelper) New() *ImageMeta {
	record := &ImageMeta{}
	// check DateCreated
	record.DateCreated = time.Now()
	return record
}

func (h *imageMetaHelper) NewFromRequest(req *http.Request) (*ImageMeta, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *imageMetaHelper) LoadAndUpdateFromRequest(req *http.Request) (*ImageMeta, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.ImageMetaID <= 0 {
		return nil, errors.New("The  failed to load because ImageMetaID was not found in the request.")
	}

	return newRecord, nil
}

func (h *imageMetaHelper) UpdateFromRequest(req *http.Request, record *ImageMeta) error {
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

func (h *imageMetaHelper) All() (ImageMetas, error) {
	var records ImageMetas
	err := h.DB.Select("*").
		From("image_meta").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *imageMetaHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (ImageMetas, error) {
	var records ImageMetas
	err := h.DB.Select("*").
		From("image_meta").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *imageMetaHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*ImageMeta, error) {
	var record ImageMeta
	if h.DB == nil {
		panic("no db")
	}
	err := h.DB.Select("*").
		From("image_meta").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (h *imageMetaHelper) Paged(pageNum int, itemsPerPage int) (*PagedData, error) {
	pd, err := h.PagedBy(pageNum, itemsPerPage, "date_created", "") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (h *imageMetaHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string, direction string) (*PagedData, error) {
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

	var records ImageMetas
	err := h.DB.Select("*").
		From("image_meta").
		OrderBy(orderByFieldName + " " + direction).
		Offset(uint64((pageNum - 1) * itemsPerPage)).
		Limit(uint64(itemsPerPage)).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	count := 0
	h.DB.SQL(`select count(image_meta_id) from image_meta`).QueryStruct(&count)
	return NewPagedData(records, orderByFieldName, direction, itemsPerPage, pageNum, count), nil
}

func (h *imageMetaHelper) Load(id int) (*ImageMeta, error) {
	record := &ImageMeta{}
	err := h.DB.
		Select("*").
		From("image_meta").
		Where("image_meta_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *imageMetaHelper) Save(record *ImageMeta) error {
	// date created always should be set, if its null just set it to now.
	if record.DateCreated.IsZero() {
		record.DateCreated = time.Now()
	}

	// was just modified
	record.DateModified = time.Now()

	err := h.save(record)
	if err != nil {
		return err
	}

	return err
}

func (h *imageMetaHelper) SaveMany(records ImageMetas) error {
	for _, record := range records {
		// date created always should be set, if its null just set it to now.
		if record.DateCreated.IsZero() {
			record.DateCreated = time.Now()
		}

		// was just modified
		record.DateModified = time.Now()
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

func (h *imageMetaHelper) save(record *ImageMeta) error {
	err := h.DB.
		Upsert("image_meta").
		Columns("name", "original", "css_height", "css_width", "container_height", "container_width", "css_left", "css_top", "date_created", "date_modified", "unique_id").
		Values(record.Name, record.Original, record.CSSHeight, record.CSSWidth, record.ContainerHeight, record.ContainerWidth, record.CSSLeft, record.CSSTop, record.DateCreated, record.DateModified, record.UniqueID).
		Where("unique_id=$1", record.UniqueID).
		Returning("image_meta_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

func (h *imageMetaHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("image_meta").
		Where("image_meta_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}
