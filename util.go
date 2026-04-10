package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	dotenv "github.com/joho/godotenv"
	dat "github.com/nerdynz/dat/dat"
	runner "github.com/nerdynz/dat/sqlx-runner"

	"github.com/jaybeecave/render"
)

func getRendererBlueprintsDir() string {
	wd, _ := loadEnv()
	tmplDir := os.Getenv("templates_dir")
	if tmplDir == "" {
		tmplDir = "./blueprints"
	}
	if strings.HasPrefix(tmplDir, "./") {
		tmplDir = wd + "/" + strings.TrimPrefix(tmplDir, "./")
	}
	return tmplDir
}

func getRenderer() *render.Render {
	r := render.New(render.Options{
		Directory: getRendererBlueprintsDir(),
		Funcs: []template.FuncMap{
			{
				"jsesc":     toJS,
				"nextIndex": nextIndex,
				"contains":  contains,
			},
		},
	})
	return r
}

func contains(s string, substr ...string) bool {
	for _, substr := range substr {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func getDBConnection() (*runner.DB, error) {
	_, err := loadEnv()
	if err != nil {
		return nil, err
	}
	//get url from ENV in the following format postgres://user:pass@192.168.8.8:5432/spaceio")
	dbURL := os.Getenv("DATABASE_URL")
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}

	username := u.User.Username()
	pass, isPassSet := u.User.Password()
	if !isPassSet {
		return nil, errors.New("no database password")
	}
	host, port, _ := net.SplitHostPort(u.Host)
	dbName := strings.Replace(u.Path, "/", "", 1)

	db, err := sql.Open("postgres", "dbname="+dbName+" user="+username+" password="+pass+" host="+host+" port="+port+" sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
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
	return runner.NewDB(db, "postgres"), nil
}

// for storing variables when running the templates
type viewBucket struct {
	Data map[string]interface{}
}

func newViewBucket() *viewBucket {
	return &viewBucket{Data: map[string]interface{}{
		"LTEqStr":     template.HTML(`<=`),
		"GTEqStr":     template.HTML(`>=`),
		"LTStr":       template.HTML(`<`),
		"GTStr":       template.HTML(`>`),
		"LTEq":        template.HTML(`<=`),
		"GTEq":        template.HTML(`>=`),
		"LT":          template.HTML(`<`),
		"GT":          template.HTML(`>`),
		"LEFT_BRACE":  template.JS(`{`),
		"RIGHT_BRACE": template.JS(`}`),
	}}
}

func toJS(s string) template.JS {
	return template.JS(s)
}

func nextIndex(i int) int {
	return i + 1
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

// findCwdAndEnv walks up from cwd looking for filename, returning the first and its cwd
func findCwdAndEnv(filename string) (string, string) {
	dir, err := os.Getwd()
	if err != nil {
		return dir, filename
	}
	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return dir, candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return dir, filename
}

func loadEnv() (string, error) {
	// Load .builder.env at startup by walking up from cwd.
	// This keeps secrets out of the MCP client config.\
	wd, envFile := findCwdAndEnv(".builder.env")
	if err := dotenv.Load(envFile); err != nil {
		return "", err
	}
	return wd + "/", nil
}

func createMigrationFile(path string, name string, direction string) (string, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	ts := fmt.Sprintf("%d", time.Now().Unix())
	filepath := filepath.Join(path, ts+"_"+name+"."+direction+".sql")
	_, err = os.Create(filepath)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
