package models

import (
	"github.com/gorilla/schema"
	runner "github.com/nerdynz/dat/sqlx-runner"
	"github.com/pinzolo/casee"
	validator "gopkg.in/go-playground/validator.v9"
	redis "github.com/go-redis/redis"
)

const NoRows = "sql: no rows in result set"

var modelValidator *validator.Validate
var modelDB *runner.DB
var modelCache *redis.Client
var modelDecoder *schema.Decoder

func Init(db *runner.DB, redis *redis.Client) {
	modelDecoder = schema.NewDecoder()
	modelDecoder.IgnoreUnknownKeys(true)
	modelDB = db
	modelCache = redis
	modelValidator = validator.New()
}

type PagedData struct {
	Sort      string      `json:"sort"`
	Direction string      `json:"direction"`
	Records   interface{} `json:"records"`
	Total     int         `json:"total"`
	PageNum   int         `json:"pageNum"`
	Limit     int         `json:"limit"`
}

func NewPagedData(records interface{}, orderBy string, direction string, itemsPerPage int, pageNum int, total int) *PagedData {
	return &PagedData{
		Records:   records,
		Direction: direction,
		Sort:      casee.ToPascalCase(orderBy),
		Limit:     itemsPerPage,
		PageNum:   pageNum,
		Total:     total,
	}
}
