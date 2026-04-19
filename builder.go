package main

import (
	"fmt"
	"go/build"
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/nerdynz/helpers"
	"github.com/pinzolo/casee"

	"bytes"

	"bufio"

	"io/ioutil"

	"os"

	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jaybeecave/render"
	"github.com/jinzhu/inflection"
	errors "github.com/kataras/go-errors"
	runner "github.com/nerdynz/dat/sqlx-runner"

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

type TableInfo struct {
	TableCatalog string `db:"table_catalog"`
	TableSchema  string `db:"table_schema"`
	TableName    string `db:"table_name"`
}

type Fields []Field

func listFields(db *runner.DB, tableName string, ignoreTsv bool) ([]*ColumnInfo, error) {
	// populate more variables from column names
	columns := []*ColumnInfo{}
	b := db.Select("column_name, data_type, is_nullable, table_name, udt_name").
		From("information_schema.columns")

	if ignoreTsv {
		b = b.Where("table_schema = $1 and table_name = $2 and column_name <> 'tsv'", "public", tableName) // field excluse
	} else {
		b = b.Where("table_schema = $1 and table_name = $2", "public", tableName) // field excluse
	}

	err := b.QueryStructs(&columns)
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func listTables(db *runner.DB) ([]*TableInfo, error) {
	ti := make([]*TableInfo, 0)

	err := db.SQL(`
SELECT t.table_catalog, t.table_schema, t.table_name
FROM information_schema.tables t
WHERE t.table_schema = 'public'
  AND t.table_name <> 'schema_migrations'
  AND NOT EXISTS (
      SELECT 1
      FROM pg_inherits i
      WHERE i.inhrelid::regclass::text = t.table_name
  )
`).QueryStructs(&ti)
	if err != nil {
		return nil, err
	}
	sort.Slice(ti, func(i, j int) bool {
		return ti[i].TableName < ti[j].TableName
	})
	return ti, nil
}

func createTable(tableName string, fields Fields, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()

	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	migPath := settingOrDefault("MIGRATION_PATH", "./migrations")
	migName := "create_" + bucket.getStr("TableName")
	upFile, err := createMigrationFile(migPath, migName, "up")
	if err != nil {
		return err
	}
	downFile, err := createMigrationFile(migPath, migName, "down")
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "create-table", upFile, bucket)
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "drop-table", downFile, bucket)
	if err != nil {
		return err
	}
	return nil
}

func createSearch(tableName string, fields Fields, r *render.Render, db *runner.DB) (string, error) {
	bucket := newViewBucket()

	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	migPath := settingOrDefault("MIGRATION_PATH", "./migrations")
	migName := "search_" + bucket.getStr("TableName")
	upFile, err := createMigrationFile(migPath, migName, "up")
	if err != nil {
		return "", err
	}
	downFile, err := createMigrationFile(migPath, migName, "down")
	if err != nil {
		return "", err
	}
	err = migrationFromTemplate(r, "create-search", upFile, bucket)
	if err != nil {
		return "", err
	}
	err = migrationFromTemplate(r, "drop-search", downFile, bucket)
	if err != nil {
		return "", err
	}
	return upFile, nil
}

func addFields(tableName string, fields []Field, r *render.Render, db *runner.DB) error {
	bucket := newViewBucket()
	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	migPath := settingOrDefault("MIGRATION_PATH", "./migrations")
	migName := "fields_" + bucket.getStr("TableName")
	upFile, err := createMigrationFile(migPath, migName, "up")
	if err != nil {
		return err
	}
	downFile, err := createMigrationFile(migPath, migName, "down")
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "add-fields", upFile, bucket)
	if err != nil {
		return err
	}
	err = migrationFromTemplate(r, "remove-fields", downFile, bucket)
	if err != nil {
		return err
	}
	return nil
}

func createBlankMigration(migrationName string, r *render.Render, db *runner.DB) error {

	// _, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"), casee.ToSnakeCase(migrationName))
	// if err != nil {
	// 	return err
	// }

	migPath := settingOrDefault("MIGRATION_PATH", "./migrations")
	migName := casee.ToSnakeCase(migrationName)
	_, err := createMigrationFile(migPath, migName, "up")
	if err != nil {
		return err
	}
	_, err = createMigrationFile(migPath, migName, "down")
	return err
}

func doMigration(r *render.Render, db *runner.DB) error {

	migPath := settingOrDefault("MIGRATION_PATH", "./migrations")

	wd, _ := loadEnv()
	err := os.Chdir(wd)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+migPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	// errs, ok := migrate.UpSync(os.Getenv("DATABASE_URL")+"?sslmode=disable", settingOrDefault("MIGRATION_PATH", "./migrations"))
	// finalError := ""
	// if ok {
	// 	// sweet
	// } else {
	// 	for _, err := range errs {
	// 		finalError += err.Error() + "\n"
	// 	}
	// }
	// if finalError != "" {
	// 	return errors.New(finalError)
	// }
	return nil
}

func createFeModel(tableName string, r *render.Render, db *runner.DB) error {
	_, err := createProtoAndTwirpBindings(tableName)
	if err != nil {
		return err
	}
	fullPath, err := createSomething(tableName, nil, r, db, "create-fe-model", settingOrDefault("SPA_API_PATH", "./spa/src/api/"), ":TableNameCamel.model.tmp.ts")
	if err != nil {
		return err
	}

	_, err = os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	return nil
}

func createProto(tableName string, r *render.Render, db *runner.DB) error {
	_, err := createSomething(tableName, nil, r, db, "create-proto", settingOrDefault("PROTO_PATH", "./proto/"), ":TableName.tmp.proto")
	return err
}

func createModel(tableName string, r *render.Render, db *runner.DB) ([]string, error) {
	createProtoAndTwirpBindings(tableName)

	paths := make([]string, 0)
	path, err := createSomething(tableName, nil, r, db, "create-model", settingOrDefault("RPC_PATH", "./rpc/:TableNameSnake/"), ":TableNameSnake.tmp.go")
	if err != nil {
		return nil, err
	}
	paths = append(paths, path)
	path, err = createSomething(tableName, nil, r, db, "create-model-helper", settingOrDefault("RPC_PATH", "./rpc/:TableNameSnake/"), ":TableNameSnake.helper.go")
	if err != nil {
		return nil, err
	}
	paths = append(paths, path)
	return paths, nil
}

// func createModel(tableName string, r *render.Render, db *runner.DB) error {
// 	return createSomething(tableName, nil, r, db, "create-model", "./rest/models", ":TableNameCamel.helper.tmp.go")
// }

func createRest(tableName string, r *render.Render, db *runner.DB) error {
	_, err := createSomething(tableName, nil, r, db, "create-rest", settingOrDefault("ACTIONS_PATH", "./rest/actions/"), ":TableNameCamelPlural.tmp.go")
	if err != nil {
		return err
	}
	return nil
}

func createProtoAndTwirpBindings(protoNameOrTableName string) (string, error) {
	wd, err := loadEnv()
	if err != nil {
		return "", err
	}
	tableName := strings.ReplaceAll(protoNameOrTableName, ".proto", "")
	tableName = helpers.SnakeCase(tableName)
	protoName := "proto/" + tableName + ".proto" // put it backz
	// goSrc := os.Getenv("GOPATH") + "/src"

	// return "", errors.New(wd)

	rpcPath := "rpc/" + helpers.SnakeCase(tableName) + "/"
	if err = os.MkdirAll(filepath.Join(wd, rpcPath), 0755); err != nil {
		return "", err
	}

	err = runCommandOrErrorInDirectory(wd, "/opt/homebrew/bin/protoc", "--proto_path", "proto", "--go_out", rpcPath, "--go_opt", "paths=source_relative", "--twirp_out", rpcPath, "--twirp_opt", "paths=source_relative", protoName)
	if err != nil {
		return "", err
	}

	// // INSERT struct tags
	resultingProto := strings.Replace("./rpc/"+helpers.SnakeCase(tableName)+"/"+helpers.SnakeCase(tableName)+".pb.go", "./", wd, -1)
	err = runCommandOrError("protoc-go-inject-tag", "-input="+resultingProto)
	if err != nil {
		return "", err
	}
	err = runCommandOrErrorInDirectory(wd+"/spa", "pnpm", "twirpscript")
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func createRPC(protoNameOrTableName string, r *render.Render, db *runner.DB) (string, error) {
	tableName, err := createProtoAndTwirpBindings(protoNameOrTableName)
	if err != nil {
		return "", err
	}

	fname, err := createSomething(tableName, nil, r, db, "create-rpc", settingOrDefault("RPC_PATH", "./rpc/:TableNameSnake/"), ":TableNameSnake.rpc.tmp.go")
	if err != nil {
		return "", err
	}
	return fname, nil
}

func createList(tableName string, r *render.Render, db *runner.DB) (string, error) {
	// err := createSomething(tableName, nil, r, db, "create-routes", settingOrDefault("SPA_ROUTE_PATH", "./spa/src/:TableNameCamel/"), ":TableNameCamel.routes.tmp.ts")
	// if err != nil {
	// 	return err
	// }

	// fullFilePath, err := createSomething(tableName, nil, r, db, "create-list", settingOrDefault("SPA_VIEW_PATH", "./spa/src/:TableNameCamel/sample/"), ":TableNamePascalList.vue")
	fullFilePath, err := createSomething(tableName, nil, r, db, "create-list", settingOrDefault("SPA_VIEW_PATH", "./tmp/"), ":TableNamePascalList.vue")
	if err != nil {
		return "", err
	}
	return fullFilePath, nil
}

func createListEdit(tableName string, r *render.Render, db *runner.DB) error {
	// err := createSomething(tableName, nil, r, db, "create-routes", settingOrDefault("SPA_ROUTE_PATH", "./spa/src/:TableNameCamel/"), ":TableNameCamel.routes.tmp.ts")
	// if err != nil {
	// 	return err
	// }
	_, err := createSomething(tableName, nil, r, db, "create-list-edit", settingOrDefault("SPA_VIEW_PATH", "./spa/src/:TableNameCamel/sample/"), ":TableNamePascalListEdit.vue")
	if err != nil {
		return err
	}
	return nil
}

func createEdit(tableName string, r *render.Render, db *runner.DB) error {
	_, err := createSomething(tableName, nil, r, db, "create-multiedit-line", settingOrDefault("SPA_VIEW_PATH", "./spa/src/:TableNameCamel/sample/"), ":TableNamePascalMultiEditLine.vue")
	if err != nil {
		return err
	}
	_, err = createSomething(tableName, nil, r, db, "create-multiedit", settingOrDefault("SPA_VIEW_PATH", "./spa/src/:TableNameCamel/sample/"), ":TableNamePascalMultiEdit.vue")
	if err != nil {
		return err
	}
	_, err = createSomething(tableName, nil, r, db, "create-edit", settingOrDefault("SPA_VIEW_PATH", "./spa/src/:TableNameCamel/sample/"), ":TableNamePascalEdit.vue")
	if err != nil {
		return err
	}
	return nil
}

func createSomething(tableName string, fields Fields, r *render.Render, db *runner.DB, tmpl string, path string, ext string) (string, error) {
	bucket := newViewBucket()
	bucket.add("TableName", tableName)
	bucket.add("Fields", fields)

	// populate variables
	// tableName := bucket.getStr("TableName")
	tableNamePascal := casee.ToPascalCase(tableName)
	tableNameCamel := casee.ToCamelCase(tableName)
	tableNameTitle := strings.Title(tableName)
	tableNameSnake := strcase.SnakeCase(tableName)
	tableNameLower := strings.ToLower(tableName)
	// tableID := tableName + "_id"
	tableULID := tableName + "_ulid"
	tnJnt := strings.Join(strings.Split(tableNamePascal, "_"), " ")

	bucket.add("TableNameSpaces", tnJnt)
	bucket.add("TableNamePascal", tableNamePascal)
	bucket.add("TableNameCamel", tableNameCamel)
	bucket.add("TableNameTitle", tableNameTitle)
	bucket.add("TableNameLower", tableNameLower)
	bucket.add("TableNamePlural", inflection.Plural(tableNameLower))
	bucket.add("TableNamePluralPascal", inflection.Plural(tableNamePascal))
	bucket.add("TableNamePluralCamel", inflection.Plural(tableNameCamel))
	bucket.add("TableNameKebab", strcase.KebabCase(tableName))
	bucket.add("TableNameSnake", strcase.SnakeCase(tableName))
	// bucket.add("TableID", tableID)
	bucket.add("TableULID", tableULID)
	bucket.add("TableULIDPascal", strings.ReplaceAll(casee.ToPascalCase(tableULID), "Ulid", "Ulid"))
	bucket.add("TableULIDCamel", strings.ReplaceAll(casee.ToCamelCase(tableULID), "Ulid", "Ulid"))
	bucket.add("TableULIDCamelWithRecord", "record."+strings.ReplaceAll(casee.ToCamelCase(tableULID), "Ulid", "Ulid"))

	bucket.add("TableUlid", tableULID)
	bucket.add("TableUlidPascal", strings.ReplaceAll(casee.ToPascalCase(tableULID), "Ulid", "Ulid"))
	bucket.add("TableUlidCamel", strings.ReplaceAll(casee.ToCamelCase(tableULID), "Ulid", "Ulid"))
	bucket.add("TableUlidCamelWithRecord", "record."+strings.ReplaceAll(casee.ToCamelCase(tableULID), "Ulid", "Ulid"))

	// populate more variables from column names
	columns := []*ColumnInfo{}
	err := db.Select("column_name, data_type, is_nullable, table_name, udt_name").
		From("information_schema.columns").
		Where("table_schema = $1 and table_name = $2 and column_name <> 'tsv'", "public", tableName). // field excluse
		QueryStructs(&columns)
	if err != nil {
		return "", err
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
				return "", err
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
	bucket.add("ColumnsStringCommaSeperated", template.HTML(colsDBConcat))
	bucket.add("ColumnsDBStrings", template.HTML(colsDBConcat))
	bucket.add("ColumnsCommaSeperatedExclusionUpdate", columnsCommaSeperatedExclusionUpdate)
	bucket.add("ColumnsUpsertConflictPairs", columnsCommaSeperatedExclusionUpdate)
	bucket.add("ColumnsCommaSeperatedPlaceholders", colsParamPlaceholders)
	bucket.add("ColumnsRecordPrefixedStrings", colsRecordPrefixedConcat)

	//child columns????
	columns = []*ColumnInfo{}

	err = db.Select("column_name, data_type, is_nullable, table_name, udt_name").
		From("information_schema.columns").
		Where("table_schema = $1 and column_name = $2 and column_name <> 'tsv' and table_name <> $3", "public", tableULID, tableName).
		QueryStructs(&columns)
	if err != nil {
		return "", err
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

	folderPath := strings.ReplaceAll(path, ":TableNamePascalPlural", inflection.Plural(tableNamePascal))
	folderPath = strings.ReplaceAll(folderPath, ":TableNameCamelPlural", inflection.Plural(tableNameCamel))
	folderPath = strings.ReplaceAll(folderPath, ":TableNamePascal", tableNamePascal)
	folderPath = strings.ReplaceAll(folderPath, ":TableNameCamel", tableNameCamel)
	folderPath = strings.ReplaceAll(folderPath, ":TableNameSnake", tableNameSnake)
	folderPath = strings.ReplaceAll(folderPath, ":TableName", tableName)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	ext = strings.ReplaceAll(ext, ":TableNamePascalPlural", inflection.Plural(tableNamePascal))
	ext = strings.ReplaceAll(ext, ":TableNameCamelPlural", inflection.Plural(tableNameCamel))
	ext = strings.ReplaceAll(ext, ":TableNamePascal", tableNamePascal)
	ext = strings.ReplaceAll(ext, ":TableNameCamel", tableNameCamel)
	ext = strings.ReplaceAll(ext, ":TableNameSnake", tableNameSnake)
	ext = strings.ReplaceAll(ext, ":TableName", tableName)
	ext = strings.ReplaceAll(ext, ":TableNameCamelULID", casee.ToCamelCase(tableULID))
	fullFilePath := folderPath + ext

	// resultingCodeFileAlreadyExists := true
	notTempFileFullPath := strings.ReplaceAll(fullFilePath, ".tmp", "")
	if _, err := os.Stat(notTempFileFullPath); os.IsNotExist(err) {
		// resulting file won't include tmp and will be ready to use as an existing file isn't already there
		// resultingCodeFileAlreadyExists = false
		fullFilePath = notTempFileFullPath
	}

	// if resultingCodeFileAlreadyExists && ignoreIfExists {
	// 	/// do nothing
	// 	return nil
	// }

	fo, err := os.Create(fullFilePath)
	if err != nil {
		return "", err
	}

	template := r.TemplateLookup(tmpl)
	if template == nil {
		return "", err
	}
	wr := bufio.NewWriter(fo)
	err = template.Execute(wr, bucket.Data)
	if err != nil {
		return "", err
	}
	wr.Flush()
	// err = ioutil.WriteFile("./migrations/"+tableName+".go", buffer.Bytes(), os.ModePerm)
	// if err != nil {
	// 	return "", err
	// }

	if err := fo.Close(); err != nil {
		return "", err
	}

	if filepath.Ext(fullFilePath) == ".go" {
		runCommandOrError("/opt/homebrew/bin/gofmt", "-s", "-w", fullFilePath)
	}

	return fullFilePath, nil
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
	if strings.Contains(strings.ToLower(colInfo.ColumnName), "ulid") {
		return "string"
	}
	if colInfo.DataType == "text" || colInfo.DataType == "character varying" || colInfo.UDTName == "citext" {
		return "string"
	}
	if colInfo.DataType == "uuid" {
		return "string"
	}
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "int32"
	}
	if colInfo.DataType == "boolean" {
		return "bool"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return "string"
	}
	return "string"
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

func migrationFromTemplate(r *render.Render, templateName string, filePath string, data *viewBucket) error {
	template := r.TemplateLookup(templateName)
	if template == nil {
		return errors.New("couldn't find the correct template")
	}
	var buffer bytes.Buffer
	wr := bufio.NewWriter(&buffer)
	err := template.Execute(wr, data)
	if err != nil {
		return err
	}
	wr.Flush()
	err = ioutil.WriteFile(filePath, buffer.Bytes(), os.ModePerm)
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

	if fi.IsDir() {
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
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// fmt.Println(path, replacement)
		newContents := ""
		if isProcFile || isPackageJSON || isDotEnv || isMainTs || isGoMod {
			newContents = strings.ReplaceAll(string(read), "Skeleton", name)
			newContents = strings.ReplaceAll(newContents, "skeleton", name)
		} else {
			newContents = strings.ReplaceAll(string(read), "github.com/nerdynz/", replacement)
		}

		if newContents != "" {
			err = os.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func replaceTextInFiles(directory, oldText, newText string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read file %s: %w", path, err)
		}

		// Replace the old text with the new text
		newContent := bytes.Replace(content, []byte(oldText), []byte(newText), -1)

		// Write updated content back to the file
		err = ioutil.WriteFile(path, newContent, info.Mode())
		if err != nil {
			return fmt.Errorf("could not write file %s: %w", path, err)
		}

		fmt.Printf("Replaced '%s' with '%s' in %s\n", oldText, newText, path)
		return nil
	})
}

func createProject(projectName string, outpath string) error {
	fullpath := build.Default.GOPATH + "/src/" + settingOrDefault("NEW_PROJECT_BASE_DIRECTORY", "github.com/nerdynz/skeleton/")
	if !strings.Contains(outpath, build.Default.GOPATH) {
		outpath = build.Default.GOPATH + "/src/" + outpath
	}
	if !strings.HasSuffix(outpath, "/") {
		outpath += "/"
	}
	// projectReplace := strings.TrimPrefix(outpath, build.Default.GOPATH+"/src/")
	if strings.Contains(projectName, "/") || outpath == build.Default.GOPATH+"/src/" || projectName == "" {
		return errors.New("Did you specify a project name and path?")
	}

	err := Copy(fullpath+".builder.env", outpath+".builder.env")
	if err != nil {
		return err
	}
	err = Copy(fullpath+"skeleton.code-workspace", outpath+projectName+".code-workspace")
	if err != nil {
		return err
	}

	folders := []string{"proto", "app", "blueprints", "rpc", "site", "tasks"}
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

	err = replaceTextInFiles(outpath, "skeleton", projectName)
	if err != nil {
		return err
	}
	return nil
}

func runCommandOrError(name string, arg ...string) error {
	return runCommandOrErrorInDirectoryRetry(nil, "", name, arg...)
}

func runCommandOrErrorInDirectory(directory string, name string, arg ...string) error {
	return runCommandOrErrorInDirectoryRetry(nil, directory, name, arg...)
}

func runCommandOrErrorInDirectoryRetry(retryErr error, directory string, name string, arg ...string) error {
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
		if retryErr != nil {
			// logrus.Error("\n" + name + " Failed to run!\n" + stderr.String())
			// logrus.Fatal(fmt.Sprint(err))
			return errors.New(stderr.String() + "afterRetryErr" + fmt.Sprint(err) + ". originalErr: " + fmt.Sprint(retryErr))
		} else {
			return runCommandOrErrorInDirectoryRetry(err, directory, strings.ReplaceAll(name, "opt/homebrew/bin", "usr/bin"), arg...)
		}
	}
	return nil
}

func settingOrDefault(key string, dflt string) string {
	wd, err := loadEnv()
	if err != nil {
		return ""
	}
	v := os.Getenv(key)
	if v == "" {
		v = dflt
	}
	return strings.Replace(v, "./", wd, -1)
}
