package models

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/gosimple/slug"
	"github.com/nerdynz/fakeTML"

	"strings"

	"sort"

	dat "github.com/nerdynz/dat"
	runner "github.com/nerdynz/dat/sqlx-runner"
	validator "gopkg.in/go-playground/validator.v9"
	redis "gopkg.in/redis.v5"
)

// Block Struct
type Block struct {
	BlockID             int            `db:"block_id" json:"BlockID"`
	Picture             string         `db:"picture" json:"Picture"`
	PictureTwo          string         `db:"picture_two" json:"PictureTwo"`
	PictureThree        string         `db:"picture_three" json:"PictureThree"`
	PictureFour         string         `db:"picture_four" json:"PictureFour"`
	PictureFive         string         `db:"picture_five" json:"PictureFive"`
	PictureSix          string         `db:"picture_six" json:"PictureSix"`
	HTML                string         `db:"html" json:"HTML"`
	HTMLTwo             string         `db:"html_two" json:"HTMLTwo"`
	HTMLThree           string         `db:"html_three" json:"HTMLThree"`
	HTMLFour            string         `db:"html_four" json:"HTMLFour"`
	HTMLFive            string         `db:"html_five" json:"HTMLFive"`
	HTMLSix             string         `db:"html_six" json:"HTMLSix"`
	PageID              int            `db:"page_id" json:"PageID"`
	DateCreated         time.Time      `db:"date_created" json:"DateCreated"`
	DateModified        time.Time      `db:"date_modified" json:"DateModified"`
	Additional          string         `db:"additional" json:"Additional"`
	AdditionalTwo       string         `db:"additional_two" json:"AdditionalTwo"`
	AdditionalThree     string         `db:"additional_three" json:"AdditionalThree"`
	AdditionalFour      string         `db:"additional_four" json:"AdditionalFour"`
	Type                string         `db:"type" json:"Type"`
	SortPosition        int            `db:"sort_position" json:"SortPosition"`
	IsDeleted           bool           `json:"IsDeleted"`
	UUID                string         `db:"uuid" json:"UUID"`
	ContentFromTable    dat.NullString `db:"content_from_table" json:"ContentFromTable"`
	ContentFromID       dat.NullInt64  `db:"content_from_id" json:"ContentFromID"`
	ContentFromTableTwo dat.NullString `db:"content_from_table_two" json:"ContentFromTableTwo"`
	ContentFromIDTwo    dat.NullInt64  `db:"content_from_id_two" json:"ContentFromIDTwo"`
	// LinkedContent       map[string]interface{}
	// LinkedContentTwo    map[string]interface{}
	LinkedContents []map[string]interface{}
}

// func (b *Block) PopulateLinkedContent(content []map[string]interface{}) {
// 	if len(content) > 0 {
// 		b.LinkedContent = make(map[string]interface{})
// 		for key, val := range content[0] {
// 			b.LinkedContent[casee.ToPascalCase(key)] = val
// 		}
// 	}
// }

// func (b *Block) PopulateLinkedContentTwo(content []map[string]interface{}) {
// 	if len(content) > 0 {
// 		b.LinkedContentTwo = make(map[string]interface{})
// 		for key, val := range content[0] {
// 			b.LinkedContentTwo[casee.ToPascalCase(key)] = val
// 		}
// 	}
// }
func (b *Block) PopulateLinkedContents(contents []map[string]interface{}) {
	b.LinkedContents = contents
}
func (b *Block) Slug() string {
	return slug.Make(b.Type)
}

var blockHelperGlobal *blockHelper

type Blocks []*Block

type blockHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func BlockHelper() *blockHelper {
	if blockHelperGlobal == nil {
		blockHelperGlobal = newBlockHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return blockHelperGlobal
}

func newBlockHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *blockHelper {
	helper := &blockHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"block_id", "picture", "picture_two", "picture_three", "picture_four", "picture_five", "picture_six", "html", "html_two", "html_three", "html_four", "html_five", "html_six", "page_id", "date_created", "date_modified", "content_from_table", "content_from_table_two", "content_from_id", "content_from_id_two", "type", "sort_position", "uuid", "additional", "additional_two", "additional_three", "additional_four"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "sort_position, date_created, date_modified"

	return helper
}

func (h *blockHelper) New() *Block {
	record := &Block{}
	// check DateCreated
	record.SortPosition = 50
	record.DateCreated = time.Now()
	return record
}

func (h *blockHelper) NewFake(pageID int) *Block {
	record := &Block{}
	// check DateCreated
	record.SortPosition = 50
	record.DateCreated = time.Now()
	record.PageID = pageID
	record.HTML = fakeTML.Combo()
	record.HTMLTwo = fakeTML.Combo()
	record.HTMLThree = fakeTML.Combo()
	record.HTMLFour = fakeTML.Combo()
	record.HTMLFive = fakeTML.Combo()
	record.HTMLSix = fakeTML.Combo()
	blkType := randomBlockType()
	record.Type = blkType
	record.Picture = randomishBlockImage(blkType, 1)
	record.PictureTwo = randomishBlockImage(blkType, 2)
	record.PictureThree = randomishBlockImage(blkType, 3)
	record.PictureFour = randomishBlockImage(blkType, 4)
	record.PictureFive = randomishBlockImage(blkType, 5)
	record.PictureSix = randomishBlockImage(blkType, 6)
	return record
}

func (h *blockHelper) KitchenBlocks(pageID int) Blocks {
	blks := make(Blocks, 0)
	for i := 0; i < 10; i++ {
		blks = append(blks, h.NewFake(pageID))
	}
	return blks
}

func (h *blockHelper) NewFromRequest(req *http.Request) (*Block, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *blockHelper) LoadAndUpdateFromRequest(req *http.Request) (*Block, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.BlockID <= 0 {
		return nil, errors.New("The  failed to load because BlockID was not found in the request.")
	}

	return newRecord, nil
}

func (h *blockHelper) UpdateFromRequest(req *http.Request, record *Block) error {
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

func (h *blockHelper) All() (Blocks, error) {
	var records Blocks
	err := h.DB.Select("*").
		From("block").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *blockHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (Blocks, error) {
	var records Blocks
	err := h.DB.Select("*").
		From("block").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *blockHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*Block, error) {
	var record Block

	err := h.DB.Select("*").
		From("block").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *blockHelper) Paged(pageNum int, itemsPerPage int) (*PagedData, error) {
	pd, err := h.PagedBy(pageNum, itemsPerPage, "date_created", "") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (h *blockHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string, direction string) (*PagedData, error) {
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

	var records Blocks
	err := h.DB.Select("*").
		From("block").
		OrderBy(orderByFieldName + " " + direction).
		Offset(uint64((pageNum - 1) * itemsPerPage)).
		Limit(uint64(itemsPerPage)).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	count := 0
	h.DB.SQL(`select count(block_id) from block`).QueryStruct(&count)
	return NewPagedData(records, orderByFieldName, direction, itemsPerPage, pageNum, count), nil
}

func (h *blockHelper) Load(id int) (*Block, error) {
	record := &Block{}
	err := h.DB.
		Select("*").
		From("block").
		Where("block_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *blockHelper) LoadByPageID(id int) (Blocks, error) {
	var records Blocks
	err := h.DB.
		Select("*").
		From("block").
		Where("page_id = $1", id).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}
	for _, block := range records {
		// CUSTOM BLOCK LOADING GOES HERE
		// _______________________________________________________________
		if block.Type == "Recent Work" || block.Type == "All Work" {
			contents := make([]map[string]interface{}, 0)
			b := h.DB.Select("*").
				From("work").
				Where("approval = ''").
				OrderBy("sort_position")
			if block.Type == "Recent Work" {
				b.Limit(2)
			}
			err := b.QueryObject(&contents)

			if err != nil {
				return nil, err
			}
			block.LinkedContents = contents
		}
	}

	return records, nil
}

func (h *blockHelper) Save(record *Block) error {
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

func (h *blockHelper) SaveMany(records Blocks) error {
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
		if record.IsDeleted {
			_, err := h.Delete(record.BlockID)
			if err != nil {
				return err
			}
		} else {
			// everything is validated so now re loop and do the actual saving... this should probably be a tx that can just rollback
			err := h.save(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *blockHelper) save(record *Block) error {
	err := h.DB.
		Upsert("block").
		Columns("picture", "picture_two", "picture_three", "picture_four", "picture_five", "picture_six", "html", "html_two", "html_three", "html_four", "html_five", "html_six", "page_id", "date_created", "date_modified", "content_from_table", "content_from_table_two", "content_from_id", "content_from_id_two", "type", "sort_position", "uuid", "additional", "additional_two", "additional_three", "additional_four").
		Values(record.Picture, record.PictureTwo, record.PictureThree, record.PictureFour, record.PictureFive, record.PictureSix, record.HTML, record.HTMLTwo, record.HTMLThree, record.HTMLFour, record.HTMLFive, record.HTMLSix, record.PageID, record.DateCreated, record.DateModified, record.ContentFromTable, record.ContentFromTableTwo, record.ContentFromID, record.ContentFromIDTwo, record.Type, record.SortPosition, record.UUID, record.Additional, record.AdditionalTwo, record.AdditionalThree, record.AdditionalFour).
		Where("block_id=$1", record.BlockID).
		Returning("block_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

// Validate a record
func (h *blockHelper) Validate(record *Block) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *blockHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("block").
		Where("block_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}

func randomBlockType() string {
	blocks := []string{
		"One Column",
		"Two Column",
		"Two Column Hero",
		"Three Column",
		"Hero Image",
		"Horizontal Line",
		"Horizontal Space",
	}
	return blocks[rand.Intn(len(blocks))]
}

func randomishBlockImage(blockType string, position int) string {
	if blockType == "Hero Image" {
		return fakeTML.Image16by9()
	}
	return ""
}
