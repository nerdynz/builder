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
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "skip",
			Usage: "skip bcomp diff",
		},
	}
	// app.Flags = flags
	app.Commands = []cli.Command{
		{
			Name:    "table",
			Aliases: []string{"t"},
			Usage:   "Create a new table [tablename] [fieldname:fieldtype]",
			Action: func(c *cli.Context) error {
				return createTable(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "fields",
			Aliases: []string{"f"},
			Usage:   "Add fields to an existing table [tablename] [fieldname:fieldtype]",
			Action: func(c *cli.Context) error {
				return addFields(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "create a model from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createModel(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "rest",
			Aliases: []string{"r"},
			Usage:   "create a restful interface from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createRest(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "create a edit page from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createEdit(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "create a list page from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createList(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "migration",
			Aliases: []string{"mi"},
			Usage:   "perform schema migration",
			Action: func(c *cli.Context) error {
				return doMigration(c, render, db)
			},
			Flags: flags,
		},
	}
	app.Run(os.Args)
}
