package main

import (
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/nerdynz/helpers"

	"bytes"

	"bufio"

	"io/ioutil"

	"os"

	"strings"

	"github.com/jaybeecave/render"
	"github.com/jinzhu/inflection"
	errors "github.com/kataras/go-errors"
	runner "github.com/nerdynz/dat/sqlx-runner"
	_ "gopkg.in/mattes/migrate.v1/driver/postgres"
	"gopkg.in/mattes/migrate.v1/file"
	"gopkg.in/mattes/migrate.v1/migrate"

	"github.com/serenize/snaker"
	"github.com/stoewer/go-strcase"
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

type Child struct {
	CamelName  string
	TableName  string
	PluralName string
}

type Field struct {
	FieldName    string
	FieldType    string
	FieldDefault string
}

type Fields []Field

func createTable(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()

	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./server/models/migrations", "create_"+bucket.getStr("TableName"))
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "create-table", file.UpFile, bucket)
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "drop-table", file.DownFile, bucket)
	if err != nil {
		return err
	}
	return nil
}

func addFields(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()
	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./server/models/migrations", "fields_"+bucket.getStr("TableName"))
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "add-fields", file.UpFile, bucket)
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "remove-fields", file.DownFile, bucket)
	if err != nil {
		return err
	}
	return nil
}

func doMigration(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	errs, ok := migrate.UpSync(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./server/models/migrations")
	finalError := ""
	if ok {
		// sweet
	} else {
		for _, err := range errs {
			finalError += err.Error() + "\n"
		}

	}
	return nil
}

func createModel(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, fields, r, db, "create-model", "./server/models/", ":TableNameCamel.go.tmp")
}

func createRest(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, fields, r, db, "create-rest", "./server/actions/", ":TableNameCamelPlural.go.tmp")
}

func createList(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	err := createSomethingNoDiff(tableName, fields, r, db, "create-list-index", "./admin/pages/:TableNameCamelPlural/", "index.vue", true)
	if err != nil {
		return err
	}
	return createSomething(tableName, fields, r, db, "create-list", "./admin/pages/:TableNameCamelPlural/", ":TableNameCamelList.vue.tmp")
}

func createEdit(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, fields, r, db, "create-edit", "./admin/pages/:TableNameCamelPlural/_ID/", ":TableNameCamelEdit.vue.tmp")
}

func createSomething(tableName string, fields Fields, r *render.Render, db *runner.DB, tmpl string, path string, ext string) error {
	return createSomethingNoDiff(tableName, fields, r, db, tmpl, path, ext, false)
}

func createSomethingNoDiff(tableName string, fields Fields, r *render.Render, db *runner.DB, tmpl string, path string, ext string, skipDiff bool) error {
	bucket := newViewBucket()
	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	// populate variables
	// tableName := bucket.getStr("TableName")
	tableNameTitle := snaker.SnakeToCamel(tableName) // this actualy gives us a TitleCase result
	tableNameCamel := snaker.SnakeToCamelLower(tableName)
	tableNameLower := strings.ToLower(tableName)
	tableID := tableName + "_id"
	tableULID := tableName + "_ulid"
	tnJnt := strings.Join(strings.Split(tableNameTitle, "_"), " ")

	bucket.add("TableNameSpaces", tnJnt)
	bucket.add("TableNameTitle", tableNameTitle)
	bucket.add("TableNameCamel", tableNameCamel)
	bucket.add("TableNameLower", tableNameLower)
	bucket.add("TableNamePlural", inflection.Plural(tableNameLower))
	bucket.add("TableNamePluralTitle", inflection.Plural(tableNameTitle))
	bucket.add("TableNamePluralCamel", inflection.Plural(tableNameCamel))
	bucket.add("TableNameKebab", strcase.KebabCase(tableName))
	bucket.add("TableID", tableID)
	bucket.add("TableULID", tableULID)
	bucket.add("TableIDTitle", snaker.SnakeToCamel(tableID))
	bucket.add("TableIDCamel", snaker.SnakeToCamelLower(snaker.SnakeToCamel(tableID)))
	bucket.add("TableIDCamelWithRecord", "record."+snaker.SnakeToCamelLower(snaker.SnakeToCamel(tableID)))

	// populate more variables from column names
	columns := []*ColumnInfo{}
	err := db.Select("column_name, data_type, is_nullable, table_name").
		From("information_schema.columns").
		Where("table_schema = $1 and table_name = $2 and column_name <> 'tsv'", "public", tableName).
		QueryStructs(&columns)
	if err != nil {
		return err
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

	columns = []*ColumnInfo{}
	err = db.Select("column_name, data_type, is_nullable, table_name").
		From("information_schema.columns").
		Where("table_schema = $1 and (column_name = $2 or column_name = $3) and column_name <> 'tsv' and table_name <> $4", "public", tableID, tableULID, tableName).
		QueryStructs(&columns)
	if err != nil {
		return err
	}
	childrenTableNames := make([]Child, 0)
	for _, col := range columns {
		colName := snaker.SnakeToCamel(col.TableName)
		colNamePlural := inflection.Plural(colName)
		childrenTableNames = append(childrenTableNames, Child{
			PluralName: colNamePlural,
			TableName:  colName,
			CamelName:  snaker.SnakeToCamelLower(colName),
		})
	}

	bucket.add("Children", childrenTableNames)

	folderPath := strings.Replace(path, ":TableNameCamelPlural", inflection.Plural(tableNameCamel), -1)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	ext = strings.Replace(ext, ":TableNameCamelPlural", inflection.Plural(tableNameCamel), -1)
	ext = strings.Replace(ext, ":TableNameCamel", tableNameCamel, -1)
	ext = strings.Replace(ext, ":TableNameCamelID", snaker.SnakeToCamelLower(tableID), -1)
	fullpath := folderPath + ext

	fo, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	template := r.TemplateLookup(tmpl)
	if template == nil {
		return err
	}
	wr := bufio.NewWriter(fo)
	err = template.Execute(wr, bucket.Data)
	if err != nil {
		return err
	}
	wr.Flush()
	// err = ioutil.WriteFile("./server/models/migrations/"+tableName+".go", buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	if err := fo.Close(); err != nil {
		return err
	}
	fullpathNoTemp := strings.Replace(fullpath, ".tmp", "", 1)
	skip := true
	// skip := c.Bool("skip") || skipDiff
	if !skip {
		diffCommand := os.Getenv("DIFF_COMMAND")
		if diffCommand == "" {
			diffCommand = "bcomp"
		}
		err = exec.Command(diffCommand, fullpath, fullpathNoTemp).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

type ColumnInfo struct {
	ColumnName string `db:"column_name"`
	TableName  string `db:"table_name"`
	DataType   string `db:"data_type"`
	IsNullable string `db:"is_nullable"`
}

func (colInfo *ColumnInfo) Label() string {
	colName := strings.Join(strings.Split(colInfo.ColumnName, "_"), " ")
	colName = strings.Title(colName)
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
	return strings.HasPrefix(colInfo.Name(), "Date")
}

func (colInfo *ColumnInfo) IsDefault() bool {
	if colInfo.IsDate() {
		return false
	}
	return true
}

func (colInfo *ColumnInfo) ColumnNameTitle() string {
	if colInfo.ColumnName == "ulid" {
		return "ULID"
	}
	s := snaker.SnakeToCamel(colInfo.ColumnName)
	s = strings.Replace(s, "Ulid", "ULID", -1)
	return s
}

func (colInfo *ColumnInfo) ColumnNameSplitTitle() string {
	return helpers.SplitTitleCase(colInfo.ColumnName)
}

func (colInfo *ColumnInfo) ColumnNameCamel() string {
	return snaker.SnakeToCamelLower(colInfo.ColumnNameTitle())
}

func (colInfo *ColumnInfo) ColumnType() string {
	// if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
	// 	return "ULID"
	// }
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" {
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

func (colInfo *ColumnInfo) IsID() bool {
	return strings.Contains(colInfo.Name(), "_id") || strings.HasSuffix(colInfo.Name(), "ID")
}

func (colInfo *ColumnInfo) IsSort() bool {
	return colInfo.Name() == "sort_position"
}

func (colInfo *ColumnInfo) ControlType() string {
	return colInfo.InputControlType()
}

func (colInfo *ColumnInfo) InputControlType() string {
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" {
		if strings.Contains(strings.ToLower(colInfo.ColumnName), "image") ||
			strings.Contains(strings.ToLower(colInfo.ColumnName), "picture") {
			return "image"
		}
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
	if template == nil {
		return errors.New("couldn't find the correct template")
	}

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

// func camelCase(str string) string {
// 	letters := strings.Split(str, "")
// 	letters[0] = strings.ToLower(letters[0])
// 	str = strings.Join(letters, "")
// 	return str
// }

func visit(path string, name string, replacement string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil
	}
	// fmt.Println(path)

	isGoFile, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	isProcFile, err := filepath.Match("Procfile", fi.Name())
	if err != nil {
		return err
	}
	isPackageJSON, err := filepath.Match("package.json", fi.Name())
	if err != nil {
		return err
	}

	isDotEnv, err := filepath.Match(".env", fi.Name())
	if err != nil {
		return err
	}

	if isGoFile || isProcFile || isPackageJSON || isDotEnv {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// fmt.Println(path, replacement)
		newContents := ""
		if isProcFile || isPackageJSON || isDotEnv {
			newContents = strings.Replace(string(read), "scaffold", name, -1)
		} else {
			newContents = strings.Replace(string(read), "github.com/nerdynz/builder/scaffold", replacement, -1)
		}

		if newContents != "" {
			err = ioutil.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func walkFiles(path string, name, replacement string) error {
	// fmt.Println("path:" + path + "    replacement:" + replacement)
	err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		return visit(path, name, replacement, fi, err)
	})
	if err != nil {
		return err
	}
	return nil
}

// func createProject(tableName string, fields Fields, r *render.Render) error {
// 	fullpath := build.Default.GOPATH + "/src/github.com/nerdynz/builder/scaffold"
// 	projectName := c.Args().First()
// 	outpath := build.Default.GOPATH + "/src/" + c.Args().Get(1)
// 	projectReplace := c.Args().Get(1)
// 	if strings.Contains(projectName, "/") || outpath == build.Default.GOPATH+"/src/" || projectName == "" {
// 		return ("Did you specify a project name and
// 	}
// 	fmt.Println("copying " + fullpath + " to " + outpath)
// 	err := Copy(fullpath, outpath)
// 	if err != nil {
// 		if strings.HasPrefix(err.Error(), "symlink") {
// 			fmt.Println(err.Error())
// 		} else {
// 			return err
// 		}
// 	}
// 	err = walkFiles(outpath, projectName, projectReplace)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
