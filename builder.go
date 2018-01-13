package main

import (
	"html/template"
	"net/http"

	"bytes"

	"bufio"

	"io/ioutil"

	"os"

	"github.com/jaybeecave/render"
	"github.com/jinzhu/inflection"
	errors "github.com/kataras/go-errors"
	_ "github.com/mattes/migrate/driver/postgres" //for migrations
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/urfave/cli"

	"strings"

	"github.com/serenize/snaker"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

type description struct {
	Name        string
	Method      string
	URL         string
	Description string
	Function    http.HandlerFunc
}

type descriptions []description

func (slice descriptions) Len() int {
	return len(slice)
}

func (slice descriptions) Less(i int, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice descriptions) Swap(i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Field struct {
	FieldName string
	FieldType string
}

type Fields []Field

func createTable(c *cli.Context, r *render.Render, db *runner.DB) error {
	// setup
	bucket := newViewBucket()
	args := c.Args()

	if !args.Present() {
		// no args
		return cli.NewExitError("ERROR: No tablename defined", 1)
	}

	// add variables for template
	bucket.addFieldDataFromContext(c)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./server/models/migrations", "create_"+bucket.getStr("TableName"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "create-table", file.UpFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "drop-table", file.DownFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func addFields(c *cli.Context, r *render.Render, db *runner.DB) error {
	// setup
	bucket := newViewBucket()
	if !c.Args().Present() {
		// no args
		return cli.NewExitError("ERROR: No tablename defined", 1)
	}

	// add variables for template
	bucket.addFieldDataFromContext(c)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./server/models/migrations", "fields_"+bucket.getStr("TableName"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "add-fields", file.UpFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "remove-fields", file.DownFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func doMigration(c *cli.Context, r *render.Render, db *runner.DB) error {
	errs, ok := migrate.UpSync(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./server/models/migrations")
	finalError := ""
	if !ok {
		for _, err := range errs {
			finalError += err.Error() + "\n"
		}
		return errors.New(finalError)
	}
	return nil
}

func createModel(c *cli.Context, r *render.Render, db *runner.DB) error {
	return createSomething(c, r, db, "create-model", "./server/models/", ".go.tmp")
}

func createRest(c *cli.Context, r *render.Render, db *runner.DB) error {
	return createSomething(c, r, db, "create-rest", "./server/actions/", ".go.tmp")
}

func createList(c *cli.Context, r *render.Render, db *runner.DB) error {
	return createSomething(c, r, db, "create-list", "./admin/pages/:TableNameCamelPlural/", "index.vue.tmp")
}

func createEdit(c *cli.Context, r *render.Render, db *runner.DB) error {
	return createSomething(c, r, db, "create-edit", "./admin/pages/:TableNameCamelPlural/", "_:TableNameCamelID.vue.tmp")
}

func createSomething(c *cli.Context, r *render.Render, db *runner.DB, tmpl string, path string, ext string) error {
	bucket := newViewBucket()
	args := c.Args()

	if !args.Present() {
		// no args
		return cli.NewExitError("ERROR: No tablename defined", 1)
	}
	// add variables for template
	bucket.addFieldDataFromContext(c)

	// populate variables
	tableName := bucket.getStr("TableName")
	tableNameTitle := snaker.SnakeToCamel(tableName)
	tableNameCamel := camelCase(tableNameTitle)
	tableNameLower := strings.ToLower(tableName)
	tableID := tableName + "_id"
	tnJnt := strings.Join(strings.Split(tableNameTitle, "_"), " ")

	bucket.add("TableNameSpaces", tnJnt)
	bucket.add("TableNameTitle", tableNameTitle)
	bucket.add("TableNameCamel", tableNameCamel)
	bucket.add("TableNameLower", tableNameLower)
	bucket.add("TableNamePlural", inflection.Plural(tableNameLower))
	bucket.add("TableNamePluralTitle", inflection.Plural(tableNameTitle))
	bucket.add("TableNamePluralCamel", inflection.Plural(tableNameCamel))
	bucket.add("TableID", tableID)
	bucket.add("TableIDTitle", snaker.SnakeToCamel(tableID))
	bucket.add("TableIDCamel", camelCase(snaker.SnakeToCamel(tableID)))

	// populate more variables from column names
	columns := []*ColumnInfo{}
	err := db.Select("column_name, data_type, is_nullable").
		From("information_schema.columns").
		Where("table_schema = $1 and table_name = $2 and column_name <> 'tsv'", "public", tableName).
		QueryStructs(&columns)
	if err != nil {
		return cli.NewExitError("error 10: "+err.Error(), 1)
	}
	colsDBConcat := `"`
	colsRecordPrefixedConcat := ""
	for i, col := range columns {
		if col.ColumnName == tableID {
			// we never want to include the table_id where these values are used because id's get generated from the database
			continue
		}
		colsDBConcat += col.ColumnName + `"`
		colsRecordPrefixedConcat += "record." + col.ColumnNameTitle()
		if i != (len(columns) - 1) {
			colsDBConcat += `, "`
			colsRecordPrefixedConcat += ", "
		}
	}
	bucket.add("Columns", columns)
	bucket.add("ColumnsDBStrings", template.HTML(colsDBConcat))
	bucket.add("ColumnsRecordPrefixedStrings", colsRecordPrefixedConcat)

	folderPath := strings.Replace(path, ":TableNameCamelPlural", inflection.Plural(tableNameCamel), -1)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	ext = strings.Replace(ext, ":TableNameCamelID", camelCase(snaker.SnakeToCamel(tableID)), -1)
	fullpath := folderPath + ext

	fo, _ := os.Create(fullpath)
	template := r.TemplateLookup(tmpl)
	wr := bufio.NewWriter(fo)
	err = template.Execute(wr, bucket.Data)
	if err != nil {
		return err
	}
	wr.Flush()
	// err = ioutil.WriteFile("./server/models/migrations/"+tableName+".go", buffer.Bytes(), os.ModePerm)
	if err != nil {
		return cli.NewExitError("error 20: "+err.Error(), 1)
	}

	if err := fo.Close(); err != nil {
		return cli.NewExitError("error 30: "+err.Error(), 1)
	}
	// exec.Command("bcomp", fullpathNoTemp, fullpath)
	return nil
}

type ColumnInfo struct {
	ColumnName string `db:"column_name"`
	DataType   string `db:"data_type"`
	IsNullable string `db:"is_nullable"`
}

func (colInfo *ColumnInfo) Label() string {
	colName := snaker.SnakeToCamel(colInfo.ColumnName)
	colName = strings.Join(strings.Split(colName, "_"), " ")
	return colName
}

func (colInfo *ColumnInfo) Name() string {
	colName := snaker.SnakeToCamel(colInfo.ColumnName)
	return colName
}

func (colInfo *ColumnInfo) IsNullField() bool {
	return colInfo.IsNullable == "YES"
}

func (colInfo *ColumnInfo) IsDate() bool {
	return strings.HasPrefix(colInfo.ColumnName, "Date")
}

func (colInfo *ColumnInfo) IsDefault() bool {
	if colInfo.IsDate() {
		return false
	}
	return true
}

func (colInfo *ColumnInfo) ColumnNameTitle() string {
	return snaker.SnakeToCamel(colInfo.ColumnName)
}

func (colInfo *ColumnInfo) ColumnNameCamel() string {
	return camelCase(colInfo.ColumnNameTitle())
}

func (colInfo *ColumnInfo) ColumnType() string {
	if colInfo.DataType == "text" {
		return "string"
	}
	if colInfo.DataType == "uuid" {
		return "string"
	}
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "int"
	}
	if colInfo.DataType == "boolean" {
		return "bool"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return "time.Time"
	}
	return ""
}

func (colInfo *ColumnInfo) InputControlType() string {
	if colInfo.DataType == "text" {
		if strings.Contains(strings.ToLower(colInfo.ColumnName), "html") {
			return "richtext"
		}
		if strings.Contains(strings.ToLower(colInfo.ColumnName), "text") {
			return "textarea"
		}
		return "text"
	}
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "number"
	}
	if colInfo.DataType == "boolean" {
		return "checkbox"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return "datetime"
	}
	return ""
}

func migrationFromTemplate(r *render.Render, templateName string, file *file.File, data *viewBucket) error {
	template := r.TemplateLookup(templateName)
	buffer := bytes.NewBuffer(file.Content)
	wr := bufio.NewWriter(buffer)
	err := template.Execute(wr, data)
	if err != nil {
		return err
	}
	wr.Flush()
	err = ioutil.WriteFile(file.Path+"/"+file.FileName, buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func camelCase(str string) string {
	letters := strings.Split(str, "")
	letters[0] = strings.ToLower(letters[0])
	str = strings.Join(letters, "")
	return str
}
