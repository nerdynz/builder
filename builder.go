package main

import (
	"fmt"
	"go/build"
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/nerdynz/helpers"
	"github.com/pinzolo/casee"
	"github.com/sirupsen/logrus"

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
	CamelName            string
	TableName            string
	PluralName           string
	PluralCamelName      string
	TableNamePluralCamel string
}

type Field struct {
	FieldName     string
	FieldType     string
	FieldDefault  string
	FieldPriority string
}

type Fields []Field

func createTable(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()

	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"), "create_"+bucket.getStr("TableName"))
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

func createSearch(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()

	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"), "search_"+bucket.getStr("TableName"))
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "create-search", file.UpFile, bucket)
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "drop-search", file.DownFile, bucket)
	if err != nil {
		return err
	}
	return nil
}

func addFields(tableName string, fields []Field, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()
	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"), "fields_"+bucket.getStr("TableName"))
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

func createBlankMigration(migrationName string, r *render.Render, db *runner.DB) error {
	_, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"), casee.ToCamelCase(migrationName))
	if err != nil {
		return err
	}
	return nil
}

func doMigration(r *render.Render, db *runner.DB) error {
	errs, ok := migrate.UpSync(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"))
	finalError := ""
	if ok {
		// sweet
	} else {
		for _, err := range errs {
			finalError += err.Error() + "\n"
		}
	}
	if finalError != "" {
		return errors.New(finalError)
	}
	return nil
}

func createAPI(tableName string, r *render.Render, db *runner.DB) error {
	runCommandOrFatalInDirectory("./spa", "pnpm", "twirpscript")
	return createSomething(tableName, nil, r, db, "create-api", settingOrDefault("SPA_API_PATH", "./spa/src/api/"), ":TableNameCamel.tmp.ts", true)
}

func createProto(tableName string, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, nil, r, db, "create-proto", settingOrDefault("PROTO_PATH", "./proto/"), ":TableName.tmp.proto", true)
}

func createModel(tableName string, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, nil, r, db, "create-model", settingOrDefault("RPC_PATH", "./rpc/:TableNameSnake/"), ":TableNameSnake.helper.tmp.go", true)
}

// func createModel(tableName string, r *render.Render, db *runner.DB) error {
// 	return createSomething(tableName, nil, r, db, "create-model", "./rest/models", ":TableNameCamel.helper.tmp.go", true)
// }

func createRest(tableName string, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, nil, r, db, "create-rest", settingOrDefault("ACTIONS_PATH", "./rest/actions/"), ":TableNameCamelPlural.tmp.go", true)
}

func createRPC(protoNameOrTableName string, r *render.Render, db *runner.DB) error {
	// if _, err := os.Stat("./proto/:TableName.proto"); os.IsNotExist(err) {
	// 	err := createProto(tableName, r, db)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// protoName := tableName + ".proto"
	tableName := strings.ReplaceAll(protoNameOrTableName, ".proto", "")
	tableName = helpers.SnakeCase(tableName)
	protoName := tableName + ".proto" // put it backz
	goSrc := os.Getenv("GOPATH") + "/src"
	runCommandOrFatal("/opt/homebrew/bin/protoc", "--proto_path", "./proto", "--go_out", goSrc, "--twirp_out", goSrc, protoName)

	// INSERT struct tags
	resultingProto := "./rpc/" + helpers.SnakeCase(tableName) + "/" + helpers.SnakeCase(tableName) + ".pb.go"
	runCommandOrFatal("protoc-go-inject-tag", "-input="+resultingProto)

	return createSomething(tableName, nil, r, db, "create-rpc", settingOrDefault("RPC_PATH", "./rpc/:TableNameSnake/"), ":TableNameSnake.rpc.tmp.go", true)
}

func createList(tableName string, r *render.Render, db *runner.DB) error {
	// err := createSomething(tableName, nil, r, db, "create-list-index", "./spa/src/views/:TableNameCamelPlural/", "index.vue", false)
	err := createSomething(tableName, nil, r, db, "create-route", settingOrDefault("SPA_ROUTE_PATH", "./spa/src/router/"), ":TableNameCamelRoute.ts", true)
	if err != nil {
		return err
	}
	return createSomething(tableName, nil, r, db, "create-list", settingOrDefault("SPA_VIEW_PATH", "./spa/src/views/:TableNameCamelPlural/"), ":TableNamePascalList.tmp.vue", true)
}

func createEdit(tableName string, r *render.Render, db *runner.DB) error {
	return createSomething(tableName, nil, r, db, "create-edit", settingOrDefault("SPA_VIEW_PATH", "./spa/src/views/:TableNameCamelPlural/"), ":TableNamePascalEdit.tmp.vue", true)
}

func createSomething(tableName string, fields Fields, r *render.Render, db *runner.DB, tmpl string, path string, ext string, diff bool) error {
	bucket := newViewBucket()
	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	// populate variables
	// tableName := bucket.getStr("TableName")
	tableNamePascal := casee.ToPascalCase(tableName)
	tableNameCamel := casee.ToCamelCase(tableName)
	tableNameSnake := strcase.SnakeCase(tableName)
	tableNameLower := strings.ToLower(tableName)
	// tableID := tableName + "_id"
	tableULID := tableName + "_ulid"
	tnJnt := strings.Join(strings.Split(tableNamePascal, "_"), " ")

	bucket.add("TableNameSpaces", tnJnt)
	bucket.add("TableNamePascal", tableNamePascal)
	bucket.add("TableNameCamel", tableNameCamel)
	bucket.add("TableNameLower", tableNameLower)
	bucket.add("TableNamePlural", inflection.Plural(tableNameLower))
	bucket.add("TableNamePluralPascal", inflection.Plural(tableNamePascal))
	bucket.add("TableNamePluralCamel", inflection.Plural(tableNameCamel))
	bucket.add("TableNameKebab", strcase.KebabCase(tableName))
	bucket.add("TableNameSnake", strcase.SnakeCase(tableName))
	// bucket.add("TableID", tableID)
	bucket.add("TableULID", tableULID)
	bucket.add("TableULIDPascal", strings.Replace(casee.ToPascalCase(tableULID), "Ulid", "Ulid", -1))
	bucket.add("TableULIDCamel", strings.Replace(casee.ToCamelCase(tableULID), "Ulid", "Ulid", -1))
	bucket.add("TableULIDCamelWithRecord", "record."+strings.Replace(casee.ToCamelCase(tableULID), "Ulid", "Ulid", -1))

	// populate more variables from column names
	columns := []*ColumnInfo{}
	err := db.Select("column_name, data_type, is_nullable, table_name, udt_name").
		From("information_schema.columns").
		Where("table_schema = $1 and table_name = $2 and column_name <> 'tsv'", "public", tableName). // field excluse
		QueryStructs(&columns)
	if err != nil {
		return err
	}
	for _, col := range columns {
		if col.DataType == "USER-DEFINED" {
			vals := make([]string, 0)
			err := db.SQL(`		
			select e.enumlabel::text as enum_value
			from pg_type t 
				 join pg_enum e on t.oid = e.enumtypid  
				 join pg_catalog.pg_namespace n ON n.oid = t.typnamespace
				where t.typname = $1
			`, col.UDTName).QuerySlice(&vals)
			if err != nil {
				return err
			}
			col.EnumValues = vals
		}
	}

	colsCommaSeperated := ``
	colsParamPlaceholders := ``
	colsDBConcat := `"`
	colsRecordPrefixedConcat := ""
	columnsCommaSeperatedExclusionUpdate := "SET "
	for i, col := range columns {
		colsParamPlaceholders += `$` + strconv.Itoa(i+1)
		colsCommaSeperated += col.ColumnName
		colsDBConcat += col.ColumnName + `"`
		colsRecordPrefixedConcat += "record." + col.ColumnNamePascal()

		if !(col.ColumnName == tableULID || col.ColumnName == "site_ulid") {
			columnsCommaSeperatedExclusionUpdate += col.ColumnName + " = EXCLUDED." + col.ColumnName
			// we never want to include the table_ulid for columnsCommaSeperatedExclusionUpdate
		}
		if i != (len(columns) - 1) {
			colsDBConcat += `, "`
			colsRecordPrefixedConcat += ", "
			colsParamPlaceholders += ", "
			colsCommaSeperated += ", "
			if !(col.ColumnName == tableULID || col.ColumnName == "site_ulid") {
				columnsCommaSeperatedExclusionUpdate += ", "
			}
		}
	}
	bucket.add("Columns", columns)
	bucket.add("ColumnsCommaSeperated", colsCommaSeperated)
	bucket.add("ColumnsCommaSeperatedExclusionUpdate", columnsCommaSeperatedExclusionUpdate)
	bucket.add("ColumnsCommaSeperatedPlaceholders", colsParamPlaceholders)
	bucket.add("ColumnsDBStrings", template.HTML(colsDBConcat))
	bucket.add("ColumnsRecordPrefixedStrings", colsRecordPrefixedConcat)

	//child columns????
	columns = []*ColumnInfo{}

	err = db.Select("column_name, data_type, is_nullable, table_name, udt_name").
		From("information_schema.columns").
		Where("table_schema = $1 and column_name = $2 and column_name <> 'tsv' and table_name <> $3", "public", tableULID, tableName).
		QueryStructs(&columns)
	if err != nil {
		return err
	}
	childrenTableNames := make([]Child, 0)
	for _, col := range columns {
		tableName := casee.ToPascalCase(col.TableName)
		tableNamePlural := inflection.Plural(tableName)
		childrenTableNames = append(childrenTableNames, Child{
			PluralName:           tableNamePlural,
			TableName:            tableName,
			CamelName:            casee.ToCamelCase(tableName),
			PluralCamelName:      inflection.Plural(tableName),
			TableNamePluralCamel: inflection.Plural(tableName),
		})
	}

	bucket.add("Children", childrenTableNames)

	folderPath := strings.Replace(path, ":TableNamePascalPlural", inflection.Plural(tableNamePascal), -1)
	folderPath = strings.Replace(folderPath, ":TableNameCamelPlural", inflection.Plural(tableNameCamel), -1)
	folderPath = strings.Replace(folderPath, ":TableNamePascal", tableNamePascal, -1)
	folderPath = strings.Replace(folderPath, ":TableNameCamel", tableNameCamel, -1)
	folderPath = strings.Replace(folderPath, ":TableNameSnake", tableNameSnake, -1)
	folderPath = strings.Replace(folderPath, ":TableName", tableName, -1)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	ext = strings.Replace(ext, ":TableNamePascalPlural", inflection.Plural(tableNamePascal), -1)
	ext = strings.Replace(ext, ":TableNameCamelPlural", inflection.Plural(tableNameCamel), -1)
	ext = strings.Replace(ext, ":TableNamePascal", tableNamePascal, -1)
	ext = strings.Replace(ext, ":TableNameCamel", tableNameCamel, -1)
	ext = strings.Replace(ext, ":TableNameSnake", tableNameSnake, -1)
	ext = strings.Replace(ext, ":TableName", tableName, -1)
	ext = strings.Replace(ext, ":TableNameCamelULID", casee.ToCamelCase(tableULID), -1)
	tempFileFullPath := folderPath + ext

	resultingCodeFileAlreadyExists := true
	notTempFileFullPath := strings.Replace(tempFileFullPath, ".tmp", "", 1)
	if _, err := os.Stat(notTempFileFullPath); os.IsNotExist(err) {
		// resulting file won't include tmp and will be ready to use as an existing file isn't already there
		resultingCodeFileAlreadyExists = false
		tempFileFullPath = notTempFileFullPath
	}

	fo, err := os.Create(tempFileFullPath)
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
	// err = ioutil.WriteFile("./migrations/"+tableName+".go", buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	if err := fo.Close(); err != nil {
		return err
	}
	if filepath.Ext(tempFileFullPath) == ".go" {
		err := exec.Command("/opt/homebrew/bin/gofmt", "-s", "-w", tempFileFullPath).Run()
		if err != nil {
			return err
		}
	}

	if resultingCodeFileAlreadyExists {
		err = os.Rename(tempFileFullPath, notTempFileFullPath+".tmp")
		if err != nil {
			return err
		}
		tempFileFullPath = notTempFileFullPath + ".tmp" // swap the tmp back to the end of the file to stop compliation errors

		// skip := os.Getenv("SKIP_DIFF") == "true" // This should be always rather than skip diff
		// if !skip && diff {
		// 	go func() {
		// 		diffCommand := os.Getenv("DIFF_COMMAND")
		// 		if diffCommand == "" {
		// 			diffCommand = "bcomp"
		// 		}
		// 		_ = exec.Command(diffCommand, tempFileFullPath, notTempFileFullPath).Run() // dont care if it errors
		// 	}()
		// }
	}
	return nil
}

type ColumnInfo struct {
	ColumnName string   `db:"column_name"`
	TableName  string   `db:"table_name"`
	DataType   string   `db:"data_type"`
	UDTName    string   `db:"udt_name"`
	IsNullable string   `db:"is_nullable"`
	EnumValues []string `db:"enum_values"`
}

func (colInfo *ColumnInfo) Label() string {
	colName := strings.Join(strings.Split(colInfo.ColumnName, "_"), " ")
	colName = strings.Title(colName)
	return colName
}

func (colInfo *ColumnInfo) Name() string {
	colName := casee.ToPascalCase(colInfo.ColumnName)
	return colName
}

func (colInfo *ColumnInfo) NameCamelCase() string {
	colName := casee.ToCamelCase(colInfo.ColumnName)
	return colName
}

func (colInfo *ColumnInfo) IsNullField() bool {
	return colInfo.IsNullable == "YES"
}

func (colInfo *ColumnInfo) IsNumber() bool {
	return colInfo.DataType == "integer" || colInfo.DataType == "numeric"
}

func (colInfo *ColumnInfo) IsDate() bool {
	return colInfo.DataType == "timestamp with time zone"
}

func (colInfo *ColumnInfo) IsDefault() bool {
	if colInfo.IsDate() {
		return false
	}
	return true
}

func (colInfo *ColumnInfo) ColumnNamePascal() string {
	if colInfo.ColumnName == "ulid" {
		return "Ulid"
	}
	s := casee.ToPascalCase(colInfo.ColumnName)
	return s
}

func (colInfo *ColumnInfo) ColumnNameSnake() string {
	if colInfo.ColumnName == "ulid" {
		return "ulid"
	}
	s := casee.ToSnakeCase(colInfo.ColumnName)
	return s
}

func (colInfo *ColumnInfo) ColumnNameSplitTitle() string {
	return helpers.SplitTitleCase(colInfo.ColumnName)
}

func (colInfo *ColumnInfo) ColumnNameCamel() string {
	if colInfo.ColumnName == "ulid" {
		return "ULID"
	}
	s := casee.ToCamelCase(colInfo.ColumnNamePascal())
	return s
}

func (colInfo *ColumnInfo) ColumnType() string { // VERY GO CENTRIC
	// if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
	// 	return "ULID"
	// }
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" || colInfo.UDTName == "citext" {
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
	if colInfo.DataType == "USER-DEFINED" {
		return "enum"
	}
	return ""
}

func (colInfo *ColumnInfo) ProtoType() string { // VERY GO CENTRIC
	// if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
	// 	return "ULID"
	// }
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" || colInfo.UDTName == "citext" {
		return "string"
	}
	if colInfo.DataType == "uuid" {
		return "string"
	}
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "int64"
	}
	if colInfo.DataType == "boolean" {
		return "bool"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return "string"
	}
	return ""
}

func (colInfo *ColumnInfo) JavascriptType() string { // VERY GO CENTRIC
	// if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
	// 	return "ULID"
	// }
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" || colInfo.UDTName == "citext" {
		return "string"
	}
	if colInfo.DataType == "uuid" {
		return "string"
	}
	if colInfo.DataType == "ulid" {
		return "string"
	}
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "number"
	}
	if colInfo.DataType == "boolean" {
		return "boolean"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return "Date"
	}
	if colInfo.DataType == "USER-DEFINED" {
		return "enum"
	}
	return ""
}

func (colInfo *ColumnInfo) JavascriptBlankValue() string { // VERY GO CENTRIC
	if strings.Contains(strings.ToLower(colInfo.ColumnName), "site_ulid") {
		return template.JSEscapeString("siteULID")
	}
	if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
		return template.JSEscapeString("ulid()")
	}
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" || colInfo.UDTName == "citext" {
		return template.JSEscapeString(`String()`)
	}
	// if colInfo.DataType == "uuid" {
	// 	return template.JSEscapeString("ulid()")
	// }
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "0"
	}
	if colInfo.DataType == "boolean" {
		return "false"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return template.JSEscapeString("new Date()")
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
		if strings.Contains(strings.ToLower(colInfo.ColumnName), "notes") {
			return "textarea"
		}
		if strings.Contains(strings.ToLower(colInfo.ColumnName), "site_ulid") {
			return ""
		}
		if strings.Contains(strings.ToLower(colInfo.ColumnName), strings.ToLower(colInfo.TableName)+"_ulid") {
			return ""
		}
		if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
			return "select"
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
	if colInfo.DataType == "USER-DEFINED" {
		return "select"
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
	isMainTs, err := filepath.Match("main.ts", fi.Name())
	if err != nil {
		return err
	}
	isGoMod, err := filepath.Match("go.mod", fi.Name())
	if err != nil {
		return err
	}

	if isGoFile || isProcFile || isPackageJSON || isDotEnv || isMainTs || isGoMod {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// fmt.Println(path, replacement)
		newContents := ""
		if isProcFile || isPackageJSON || isDotEnv || isMainTs || isGoMod {
			newContents = strings.Replace(string(read), "Skeleton", name, -1)
			newContents = strings.Replace(newContents, "skeleton", name, -1)
		} else {
			newContents = strings.Replace(string(read), "github.com/nerdynz/skeleton/", replacement, -1)
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

func replaceNameInFiles(path, name, replacement string) error {
	// fmt.Println("path:" + path + "    replacement:" + replacement)
	err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		return visit(path, name, replacement, fi, err)
	})
	if err != nil {
		return err
	}
	return nil
}

func createProject(projectName string, outpath string) error {
	fullpath := build.Default.GOPATH + "/src/" + settingOrDefault("NEW_PROJECT_BASE_DIRECTORY", "github.com/nerdynz/skeleton/")
	if !strings.Contains(outpath, build.Default.GOPATH) {
		outpath = build.Default.GOPATH + "/src/" + outpath
	}
	if !strings.HasSuffix(outpath, "/") {
		outpath += "/"
	}
	projectReplace := strings.TrimPrefix(outpath, build.Default.GOPATH+"/src/")
	if strings.Contains(projectName, "/") || outpath == build.Default.GOPATH+"/src/" || projectName == "" {
		return errors.New("Did you specify a project name and path?")
	}

	err := Copy(fullpath+"skeleton.code-workspace", outpath+projectName+".code-workspace")
	if err != nil {
		return err
	}

	folders := []string{"proto", "spa", "blueprints", "rpc"}
	for _, folder := range folders {
		// fmt.Println("copying " + fullpath + " to " + outpath)
		err := Copy(fullpath+folder, outpath+folder)
		if err != nil {
			if strings.HasPrefix(err.Error(), "symlink") {
				fmt.Println(err.Error())
			} else {
				return err
			}
		}
	}

	err = replaceNameInFiles(outpath, projectName, projectReplace)
	if err != nil {
		return err
	}
	return nil
}

func runCommandOrFatal(name string, arg ...string) {
	runCommandOrFatalInDirectory("", name, arg...)
}

func runCommandOrFatalInDirectory(directory string, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	if directory != "" {
		cmd.Dir = directory
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		logrus.Error("\n" + name + " Failed to run!\n" + stderr.String())
		logrus.Fatal(fmt.Sprint(err))
	}
}

func settingOrDefault(key string, dflt string) string {
	v := os.Getenv(key)
	if v == "" {
		return dflt
	}
	return v
}
