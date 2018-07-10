package main

import (
	"os"

	dotenv "github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	err := dotenv.Load() // load from .env file where scaffold is run
	if err != nil {
		panic(err)
	}
	render := getRenderer()
	db := getDBConnection()
	app := cli.NewApp()
	app.Name = "builder"
	app.Usage = "code generation"
	app.Commands = []cli.Command{
		{
			Name:    "table",
			Aliases: []string{"t"},
			Usage:   "Create a new table [tablename] [fieldname:fieldtype]",
			Action: func(c *cli.Context) error {
				return createTable(c, render, db)
			},
		},
		{
			Name:    "fields",
			Aliases: []string{"f"},
			Usage:   "Add fields to an existing table [tablename] [fieldname:fieldtype]",
			Action: func(c *cli.Context) error {
				return addFields(c, render, db)
			},
		},
		{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "create a model from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createModel(c, render, db)
			},
		},
		{
			Name:    "rest",
			Aliases: []string{"r"},
			Usage:   "create a restful interface from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createRest(c, render, db)
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "create a edit page from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createEdit(c, render, db)
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "create a list page from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createList(c, render, db)
			},
		},
		{
			Name:    "migration",
			Aliases: []string{"mi"},
			Usage:   "perform schema migration",
			Action: func(c *cli.Context) error {
				return doMigration(c, render, db)
			},
		},
	}
	app.Run(os.Args)
}
