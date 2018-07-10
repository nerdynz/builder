package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	dat "github.com/nerdynz/dat"
	runner "github.com/nerdynz/dat/sqlx-runner"

	"github.com/jaybeecave/render"
	"github.com/urfave/cli"
)

func getRenderer() *render.Render {
	r := render.New(render.Options{
		Directory: "./blueprints",
		Funcs: []template.FuncMap{
			template.FuncMap{
				"jsesc": toJS,
			},
		},
	})
	return r
}

func getDBConnection() *runner.DB {
	//get url from ENV in the following format postgres://user:pass@192.168.8.8:5432/spaceio")
	dbURL := os.Getenv("DATABASE_URL")
	u, err := url.Parse(dbURL)
	if err != nil {
		panic(err)
	}

	username := u.User.Username()
	pass, isPassSet := u.User.Password()
	if !isPassSet {
		panic("no database password")
	}
	host, port, _ := net.SplitHostPort(u.Host)
	dbName := strings.Replace(u.Path, "/", "", 1)

	db, _ := sql.Open("postgres", "dbname="+dbName+" user="+username+" password="+pass+" host="+host+" port="+port+" sslmode=disable")
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	// db connection
	return runner.NewDB(db, "postgres")
}

// for storing variables when running the templates
type viewBucket struct {
	Data map[string]interface{}
}

func newViewBucket() *viewBucket {
	return &viewBucket{Data: map[string]interface{}{
		"LTEqStr": template.HTML(`<=`),
		"GTEqStr": template.HTML(`>=`),
		"LTStr":   template.HTML(`<`),
		"GTStr":   template.HTML(`>`),
	}}
}

func toJS(s string) template.JS {
	return template.JS(s)
}

func (viewBucket *viewBucket) add(key string, value interface{}) {
	viewBucket.Data[key] = value
}

func (viewBucket *viewBucket) getStrSafe(key string) (string, error) {
	val := viewBucket.Data[key]
	if val == nil {
		return "", errors.New("could not find " + key)
	}
	strVal, ok := val.(string)
	if !ok {
		return "", errors.New("could not cast " + key + " to string")
	}
	return strVal, nil
}

// getStr - returns a string for the provided key. Will panic if key not found
func (viewBucket *viewBucket) getStr(key string) string {
	val, err := viewBucket.getStrSafe(key)
	if err != nil {
		panic(err)
	}
	return val
}

func (viewBucket *viewBucket) addFieldDataFromContext(c *cli.Context) {
	args := c.Args()
	viewBucket.add("TableName", args.First())

	fields := Fields{}
	for _, arg := range args {
		if args.First() == arg {
			continue // we dont care about the first arg as its the TableName
		}
		if strings.Contains(arg, ":") {
			strSlice := strings.Split(arg, ":")
			field := Field{
				FieldName: strSlice[0],
				FieldType: strSlice[1],
			}
			fields = append(fields, field)
		}
	}
	viewBucket.add("Fields", fields)
}
