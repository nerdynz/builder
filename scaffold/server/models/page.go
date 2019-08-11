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

// Page Struct
type Page struct {
	PageID         int       `db:"page_id" json:"PageID"`
	Title          string    `db:"title" json:"Title"`
	Subtitle       string    `db:"subtitle" json:"Subtitle"`
	Slug           string    `db:"slug" json:"Slug"`
	Summary        string    `db:"summary" json:"Summary"`
	Keywords       string    `db:"keywords" json:"Keywords"`
	PreviewPicture string    `db:"preview_picture" json:"PreviewPicture"`
	Kind           string    `db:"kind" json:"Kind"`
	HTML           string    `db:"html" json:"HTML"`
	HTMLTwo        string    `db:"html_two" json:"HTMLTwo"`
	HTMLThree      string    `db:"html_three" json:"HTMLThree"`
	HTMLFour       string    `db:"html_four" json:"HTMLFour"`
	HTMLFive       string    `db:"html_five" json:"HTMLFive"`
	HTMLSix        string    `db:"html_six" json:"HTMLSix"`
	HTMLSeven      string    `db:"html_seven" json:"HTMLSeven"`
	HTMLEight      string    `db:"html_eight" json:"HTMLEight"`
	HTMLNine       string    `db:"html_nine" json:"HTMLNine"`
	HTMLTen        string    `db:"html_ten" json:"HTMLTen"`
	HTMLEleven     string    `db:"html_eleven" json:"HTMLEleven"`
	HTMLTwelve     string    `db:"html_twelve" json:"HTMLTwelve"`
	Picture        string    `db:"picture" json:"Picture"`
	PictureTwo     string    `db:"picture_two" json:"PictureTwo"`
	PictureThree   string    `db:"picture_three" json:"PictureThree"`
	PictureFour    string    `db:"picture_four" json:"PictureFour"`
	PictureFive    string    `db:"picture_five" json:"PictureFive"`
	PictureSix     string    `db:"picture_six" json:"PictureSix"`
	Misc           string    `db:"misc" json:"Misc"`
	MiscTwo        string    `db:"misc_two" json:"MiscTwo"`
	MiscThree      string    `db:"misc_three" json:"MiscThree"`
	MiscFour       string    `db:"misc_four" json:"MiscFour"`
	MiscFive       string    `db:"misc_five" json:"MiscFive"`
	MiscSix        string    `db:"misc_six" json:"MiscSix"`
	ShowTitle      bool      `db:"show_title" json:"ShowTitle"`
	ShowSubtitle   bool      `db:"show_subtitle" json:"ShowSubtitle"`
	IsLockedSlug   bool      `db:"is_locked_slug" json:"IsLockedSlug"`
	IsSpecialPage  bool      `db:"is_special_page" json:"IsSpecialPage"`
	SpecialPageFor string    `db:"special_page_for" json:"SpecialPageFor"`
	DateCreated    time.Time `db:"date_created" json:"DateCreated"`
	DateModified   time.Time `db:"date_modified" json:"DateModified"`
	ShowInNav      string    `db:"show_in_nav" json:"ShowInNav"`
	SortPosition   int       `db:"sort_position" json:"SortPosition"`
	SEOTitle       string    `db:"seo_title" json:"SEOTitle"`
	UUID           string    `db:"uuid" json:"UUID"`
	Blocks         Blocks    `json:"Blocks"`
	IsBeingEdited  bool
	Color          string `db:"color" json:"Color"`
}

func (page *Page) HasPictures() bool {
	if page.Picture != "" {
		return true
	}
	if page.PictureTwo != "" {
		return true
	}
	if page.PictureThree != "" {
		return true
	}
	if page.PictureFour != "" {
		return true
	}
	if page.PictureFive != "" {
		return true
	}
	if page.PictureSix != "" {
		return true
	}
	return false
}

func (page *Page) LoadBlocks() error {
	blocks, err := BlockHelper().LoadByPageID(page.PageID)
	if err != nil {
		return err
	}
	page.Blocks = blocks
	return nil
}

var pageHelperGlobal *pageHelper

type Pages []*Page

type pageHelper struct {
	DB            *runner.DB
	Cache         *redis.Client
	Validator     *validator.Validate
	structDecoder *schema.Decoder
	fieldNames    []string
	orderBy       string
}

func PageHelper() *pageHelper {
	if pageHelperGlobal == nil {
		pageHelperGlobal = newPageHelper(modelDB, modelCache, modelValidator, modelDecoder)
	}
	return pageHelperGlobal
}

func newPageHelper(db *runner.DB, redis *redis.Client, validate *validator.Validate, structDecoder *schema.Decoder) *pageHelper {
	helper := &pageHelper{}
	helper.DB = db
	helper.Cache = redis
	helper.Validator = validate
	helper.structDecoder = structDecoder

	// Fields
	fieldnames := []string{"page_id", "title", "slug", "summary", "keywords", "preview_picture", "kind", "html", "html_two", "html_three", "html_four", "html_five", "html_six", "html_seven", "html_eight", "html_nine", "html_ten", "html_eleven", "html_twelve", "picture", "picture_two", "picture_three", "picture_four", "picture_five", "picture_six", "misc", "misc_two", "misc_three", "misc_four", "misc_five", "misc_six", "is_locked_slug", "is_special_page", "special_page_for", "date_created", "date_modified", "show_in_nav", "subtitle", "sort_position", "uuid", "show_title", "show_subtitle", "seo_title", "color"}
	sort.Strings(fieldnames) // sort it makes searching it work correctly
	helper.fieldNames = fieldnames
	helper.orderBy = "sort_position, date_created, date_modified"

	return helper
}

func (h *pageHelper) New() *Page {
	record := &Page{}
	// check DateCreated
	record.DateCreated = time.Now()
	record.Kind = "Standard"
	record.ShowInNav = "Top Nav"
	record.SortPosition = 500
	return record
}

func (h *pageHelper) NewFromRequest(req *http.Request) (*Page, error) {
	record := h.New()
	err := h.UpdateFromRequest(req, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *pageHelper) LoadAndUpdateFromRequest(req *http.Request) (*Page, error) {
	// dummyPack is used to get the pack ID from the request and also to check the date modified
	newRecord, err := h.NewFromRequest(req)
	if err != nil {
		return nil, err
	}

	if newRecord.PageID <= 0 {
		return nil, errors.New("The  failed to load because PageID was not found in the request.")
	}

	return newRecord, nil
}

func (h *pageHelper) UpdateFromRequest(req *http.Request, record *Page) error {
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

func (h *pageHelper) All() (Pages, error) {
	var records Pages
	err := h.DB.Select("*").
		From("page").
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *pageHelper) Where(whereSQLOrMap interface{}, args ...interface{}) (Pages, error) {
	var records Pages
	err := h.DB.Select("*").
		From("page").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (h *pageHelper) One(whereSQLOrMap interface{}, args ...interface{}) (*Page, error) {
	var record Page

	err := h.DB.Select("*").
		From("page").
		Where(whereSQLOrMap, args...).
		OrderBy(h.orderBy).
		Limit(1).
		QueryStruct(&record)

	if err != nil {
		return nil, err
	}

	err = record.LoadBlocks()
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (h *pageHelper) Paged(pageNum int, itemsPerPage int) (*PagedData, error) {
	pd, err := h.PagedBy(pageNum, itemsPerPage, "date_created", "") // date_created should be the most consistant because it doesn't change
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (h *pageHelper) PagedBy(pageNum int, itemsPerPage int, orderByFieldName string, direction string) (*PagedData, error) {
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

	var records Pages
	err := h.DB.Select("*").
		From("page").
		OrderBy(orderByFieldName + " " + direction).
		Offset(uint64((pageNum - 1) * itemsPerPage)).
		Limit(uint64(itemsPerPage)).
		QueryStructs(&records)

	if err != nil {
		return nil, err
	}

	count := 0
	h.DB.SQL(`select count(page_id) from page`).QueryStruct(&count)
	return NewPagedData(records, orderByFieldName, direction, itemsPerPage, pageNum, count), nil
}

func (h *pageHelper) Load(id int) (*Page, error) {
	record := &Page{}
	err := h.DB.
		Select("*").
		From("page").
		Where("page_id = $1", id).
		QueryStruct(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (h *pageHelper) LoadBySlug(slug string) (*Page, error) {
	record, err := h.One("slug = $1", slug)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *pageHelper) LoadBySpecialPage(slug string) (*Page, error) {
	record, err := h.One("special_page_for = $1", slug)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *pageHelper) Save(record *Page) error {
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

	for _, block := range record.Blocks {
		block.PageID = record.PageID
	}
	err = BlockHelper().SaveMany(record.Blocks)
	return err
}

func (h *pageHelper) SaveMany(records Pages) error {
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

func (h *pageHelper) save(record *Page) error {
	err := h.DB.
		Upsert("page").
		Columns("title", "slug", "summary", "keywords", "preview_picture", "kind", "html", "html_two", "html_three", "html_four", "html_five", "html_six", "html_seven", "html_eight", "html_nine", "html_ten", "html_eleven", "html_twelve", "picture", "picture_two", "picture_three", "picture_four", "picture_five", "picture_six", "misc", "misc_two", "misc_three", "misc_four", "misc_five", "misc_six", "is_locked_slug", "is_special_page", "special_page_for", "date_created", "date_modified", "show_in_nav", "subtitle", "sort_position", "uuid", "show_title", "show_subtitle", "seo_title", "color").
		Values(record.Title, record.Slug, record.Summary, record.Keywords, record.PreviewPicture, record.Kind, record.HTML, record.HTMLTwo, record.HTMLThree, record.HTMLFour, record.HTMLFive, record.HTMLSix, record.HTMLSeven, record.HTMLEight, record.HTMLNine, record.HTMLTen, record.HTMLEleven, record.HTMLTwelve, record.Picture, record.PictureTwo, record.PictureThree, record.PictureFour, record.PictureFive, record.PictureSix, record.Misc, record.MiscTwo, record.MiscThree, record.MiscFour, record.MiscFive, record.MiscSix, record.IsLockedSlug, record.IsSpecialPage, record.SpecialPageFor, record.DateCreated, record.DateModified, record.ShowInNav, record.Subtitle, record.SortPosition, record.UUID, record.ShowTitle, record.ShowSubtitle, record.SEOTitle, record.Color).
		Where("page_id=$1", record.PageID).
		Returning("page_id").
		QueryStruct(record)

	if err != nil {
		return err
	}

	return nil
}

// Validate a record
func (h *pageHelper) Validate(record *Page) (bool, error) {
	validationErrors := h.Validator.Struct(record)
	if validationErrors != nil {
		return false, validationErrors
	}
	return true, nil
}

func (h *pageHelper) Delete(recordID int) (bool, error) {
	result, err := h.DB.
		DeleteFrom("page").
		Where("page_id=$1", recordID).
		Exec()

	if err != nil {
		return false, err
	}

	return (result.RowsAffected > 0), nil
}

type NavItem struct {
	Title string `db:"title"`
	Slug  string `db:"slug"`
}

type NavItems []*NavItem

func (ni *NavItem) URL() {

}

func (h *pageHelper) LoadTopNav() (NavItems, error) {
	var navItems NavItems
	err := h.DB.SQL(`
	select title, coalesce(slug, special_page_for) as slug from page
	where show_in_nav = $1
	order by sort_position
	`, "Top Nav").QueryStructs(&navItems)
	if err != nil {
		return nil, err
	}
	return navItems, nil
}
func (h *pageHelper) LoadSideNav() (NavItems, error) {
	var navItems NavItems
	err := h.DB.SQL(`
	select title, coalesce(slug, special_page_for) as slug from page
	where show_in_nav = $1
	order by sort_position
	`, "Side Nav").QueryStructs(&navItems)
	if err != nil {
		return nil, err
	}
	return navItems, nil
}

func (h *pageHelper) LoadFooterNav() (NavItems, error) {
	var navItems NavItems
	err := h.DB.SQL(`
	select title, coalesce(slug, special_page_for) as slug from page
	where show_in_nav = $1
	order by sort_position
	`, "Footer Nav").QueryStructs(&navItems)
	if err != nil {
		return nil, err
	}
	return navItems, nil
}

func (h *pageHelper) LoadKitchenSink() (*Page, error) {
	page := h.New()
	page.Title = "A Pretty Decent Title"
	page.Subtitle = "A Pretty Decent Title"
	page.SpecialPageFor = "kitchen-sink"
	page.Blocks = BlockHelper().KitchenBlocks(page.PageID)
	return page, nil
}
