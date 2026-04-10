package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

func main() {

	s := server.NewMCPServer("builder", "1.0.0")

	s.AddTool(mcp.NewTool("scaffold_project",
		mcp.WithDescription("Scaffold a new project from the skeleton template"),
		mcp.WithString("project_name", mcp.Required(), mcp.Description("Project name")),
		mcp.WithString("project_path", mcp.Required(), mcp.Description("Output path e.g. github.com/myorg/myproject")),
	), scaffoldProject)

	s.AddTool(mcp.NewTool("create_migration",
		mcp.WithDescription("Create a blank database migration file"),
		mcp.WithString("description", mcp.Required(), mcp.Description("Short description of the migration")),
	), createMigration)

	s.AddTool(mcp.NewTool("create_table",
		mcp.WithDescription("Create a new database table with an up/down migration"),
		mcp.WithString("table", mcp.Required(), mcp.Description("Table name in snake_case")),
		mcp.WithString("fields", mcp.Required(), mcp.Description(`JSON array: [{"FieldName":"col","FieldType":"text","FieldDefault":"","FieldPriority":""}]`)),
	), createTableTool)

	s.AddTool(mcp.NewTool("list_tables"), listTablesTool)

	s.AddTool(mcp.NewTool("list_fields",
		mcp.WithDescription("list fields for a specified table"),
		mcp.WithString("table", mcp.Required(), mcp.Description("Table name in snake_case")),
	), listFieldsTool)

	s.AddTool(mcp.NewTool("add_fields",
		mcp.WithDescription("Add columns to existing tables via a migration"),
		mcp.WithString("table", mcp.Required(), mcp.Description("table name in snake_case")),
		mcp.WithString("fields", mcp.Required(), mcp.Description(`JSON array: [{"FieldName":"col","FieldType":"text","FieldDefault":"","FieldPriority":""}]`)),
	), addFieldsTool)

	s.AddTool(mcp.NewTool("generate_model",
		mcp.WithDescription("Generate Go model files for tables"),
		mcp.WithString("table", mcp.Required(), mcp.Description("table name in snake_case")),
	), generateModel)

	s.AddTool(mcp.NewTool("generate_proto",
		mcp.WithDescription("Generate .proto files for tables"),
		mcp.WithString("table", mcp.Required(), mcp.Description("table name in snake_case")),
	), generateProto)

	s.AddTool(mcp.NewTool("generate_twirp",
		mcp.WithDescription("Generate twirp bindings for  provided proto file. This action is idempotent."),
		mcp.WithString("proto", mcp.Required(), mcp.Description("name of proto file")),
	), generateTwirp)

	s.AddTool(mcp.NewTool("generate_rpc",
		mcp.WithDescription("Generate rpc server for provided proto file"),
		mcp.WithString("proto", mcp.Required(), mcp.Description("name of proto file")),
	), generateRPC)

	s.AddTool(mcp.NewTool("generate_api",
		mcp.WithDescription("Generate frontend API model (TypeScript) files for tables"),
		mcp.WithString("table", mcp.Required(), mcp.Description("table name in snake_case")),
	), generateAPI)

	s.AddTool(mcp.NewTool("generate_edit",
		mcp.WithDescription("Generate frontend Vue edit component for tables"),
		mcp.WithString("table", mcp.Required(), mcp.Description("table name in snake_case")),
	), generateEdit)

	s.AddTool(mcp.NewTool("generate_list",
		mcp.WithDescription("Generate frontend Vue list component for the provided table"),
		mcp.WithString("table", mcp.Required(), mcp.Description("table name")),
	), generateList)

	// s.AddTool(mcp.NewTool("generate_search",
	// 	mcp.WithDescription("Generate full-text search migration for a table"),
	// 	mcp.WithString("table", mcp.Required(), mcp.Description("Table name")),
	// 	mcp.WithString("fields", mcp.Required(), mcp.Description(`JSON array: [{"FieldName":"col","FieldPriority":"a"}]`)),
	// ), generateSearch)

	s.AddTool(mcp.NewTool("migrate",
		mcp.WithDescription("Run pending database migrations"),
	), runMigrate)

	s.AddTool(mcp.NewTool("ping",
		mcp.WithDescription("Echo back the resolved DATABASE_URL and working directory — useful for verifying config"),
	), pong)

	logrus.Info("starting builder")
	if err := server.ServeStdio(s); err != nil {
		logrus.Fatal(err)
	}
}

func splitTables(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func parseFields(req mcp.CallToolRequest) (Fields, error) {
	fieldsJSON, err := req.RequireString("fields")
	if err != nil {
		return nil, err
	}
	var fields Fields
	if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
		return nil, fmt.Errorf("invalid fields JSON: %w", err)
	}
	return fields, nil
}

func scaffoldProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("project_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	path, err := req.RequireString("project_path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := createProject(name, path); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText("project scaffolded: " + name), nil
}

func createMigration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	desc, err := req.RequireString("description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := createBlankMigration(desc, getRenderer(), db); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText("migration created: " + desc), nil
}

func createTableTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fields, err := parseFields(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := createTable(tableName, fields, getRenderer(), db); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("table %s created", tableName)), nil
}

func listTablesTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	tables, err := listTables(db)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var sb strings.Builder
	for _, t := range tables {
		sb.WriteString("• " + t.TableName + "\n")
	}
	return mcp.NewToolResultText(sb.String()), nil
}

func listFieldsTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fields, err := listFields(db, tableName, false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var sb strings.Builder
	for _, t := range fields {
		sb.WriteString("• " + t.ColumnName + "(" + t.DataType + ")\n")
	}
	return mcp.NewToolResultText(sb.String()), nil
}

func addFieldsTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	table, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fields, err := parseFields(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	r := getRenderer()

	if err := addFields(table, fields, r, db); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("fields added to: " + table), nil
}

func generateModel(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	table, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	r := getRenderer()
	paths, err := createModel(table, r, db)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var sb strings.Builder
	sb.WriteString("model generated for: " + table + " with the following files generated: \n")
	for _, p := range paths {
		sb.WriteString(p + "\n")
	}
	sb.WriteString("AI please check the file for any TODO: AI and do your best to respect the comments")
	return mcp.NewToolResultText(sb.String()), nil
}

func generateProto(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	table, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	r := getRenderer()
	if err := createProto(table, r, db); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText("protos generated: " + table), nil
}

func generateTwirp(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	protoName, err := req.RequireString("proto")
	if err != nil {
		return mcp.NewToolResultError("name: " + err.Error()), nil
	}
	_, err = createProtoAndTwirpBindings(protoName)
	return mcp.NewToolResultText("twirp bindings generated for " + protoName), nil
}

func generateRPC(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	protoName, err := req.RequireString("proto")
	if err != nil {
		return mcp.NewToolResultError("name: " + err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError("db: " + err.Error()), nil
	}
	// return mcp.NewToolResultText(protoName), nil
	r := getRenderer()
	fname, err := createRPC(protoName, r, db)
	if err != nil {
		return mcp.NewToolResultError("render: " + err.Error()), nil
	}
	// // b, err := os.ReadFile(fname)
	// // if err != nil {
	// // 	return mcp.NewToolResultError(err.Error()), nil
	// // }
	return mcp.NewToolResultText(filepath.Base(fname)), nil
	// return mcp.NewToolResultText(filepath.Base(fname) + "\n" + string(b)), nil
}

func generateAPI(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	table, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	r := getRenderer()
	if err := createFeModel(table, r, db); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText("API models generated: " + table), nil
}

func generateEdit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	table, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	r := getRenderer()
	if err := createEdit(table, r, db); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText("edit components generated: " + table), nil
}

func generateList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError("no table provided. use list_tables to see available tables"), nil
	}

	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	r := getRenderer()
	// sb.WriteString("The following list pages were created:\n")
	fname, err := createList(tableName, r, db)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	b, err := os.ReadFile(fname)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(filepath.Base(fname) + "\n" + string(b)), nil
}

func generateSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := req.RequireString("table")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fields, err := parseFields(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	searchFilePath, err := createSearch(tableName, fields, getRenderer(), db)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	file, err := os.ReadFile(searchFilePath)
	if err != nil {
		return mcp.NewToolResultError("file read failed: " + err.Error()), nil
	}

	return mcp.NewToolResultText("here is the output file at path: " + searchFilePath + " \n\n```" + string(file) + "```\n\n NOTE To AI please check the final syntax as there might be some errors"), nil
}

func pong(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if _, err := loadEnv(); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if os.Getenv("DATABASE_URL") == "" {
		return mcp.NewToolResultError("DATABASE_URL is not avaliable"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Working Dir: %s\nDATABASE_URL: %s", wd, redactURL(os.Getenv("DATABASE_URL")))), nil
}

func redactURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.User == nil {
		return raw
	}
	u.User = url.UserPassword(u.User.Username(), "REDACTED")
	return u.String()
}

func runMigrate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	db, err := getDBConnection()
	if err != nil {
		return mcp.NewToolResultError("connection failed " + err.Error()), nil
	}
	if err := doMigration(getRenderer(), db); err != nil {
		return mcp.NewToolResultError("migrate error " + err.Error()), nil
	}
	return mcp.NewToolResultText("migrations applied"), nil
}
