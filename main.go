package main

import (
	"os"

	dotenv "github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	render := getRenderer()
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
			Name:    "scaffold",
			Aliases: []string{"s"},
			Usage:   "Create a new project [projectname] [project path]",
			Action: func(c *cli.Context) error {
				return createProject(c, render)
			},
			Flags: flags,
		},
		{
			Name:    "table",
			Aliases: []string{"t"},
			Usage:   "Create a new table [tablename] [fieldname:fieldtype:fielddefault]",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					return err
				}
				db := getDBConnection()
				return createTable(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "fields",
			Aliases: []string{"f"},
			Usage:   "Add fields to an existing table [tablename] [fieldname:fieldtype:fielddefault]",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					return err
				}
				db := getDBConnection()
				return addFields(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "create a model from a table [tablename]",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					return err
				}
				db := getDBConnection()
				return createModel(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "rest",
			Aliases: []string{"r"},
			Usage:   "create a restful interface from a table [tablename]",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					return err
				}
				db := getDBConnection()
				return createRest(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "create a edit page from a table [tablename]",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					return err
				}
				db := getDBConnection()
				return createEdit(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "create a list page from a table [tablename]",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					return err
				}
				db := getDBConnection()
				return createList(c, render, db)
			},
			Flags: flags,
		},
		{
			Name:    "migration",
			Aliases: []string{"mi"},
			Usage:   "perform schema migration",
			Action: func(c *cli.Context) error {
				err := dotenv.Load() // load from .env file where scaffold is run
				if err != nil {
					panic(err)
				}
				db := getDBConnection()
				err = doMigration(c, render, db)
				if err != nil {
					panic(err)
				}
				return err
			},
			Flags: flags,
		},
	}
	app.Run(os.Args)
}
