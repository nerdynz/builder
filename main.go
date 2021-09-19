package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dotenv "github.com/joho/godotenv"
)

const BREAK = "\n"
const DOUBLE_BREAK = "\n\n"

var headingStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	MarginTop(1).
	PaddingTop(1).
	PaddingLeft(1).
	PaddingBottom(1).
	Width(50)

var selectionStyle = lipgloss.NewStyle().
	Bold(false).
	Foreground(lipgloss.Color("#1d313b")).
	Background(lipgloss.Color("#ffdb00")).
	// PaddingTop(2).
	// PaddingLeft(4).
	Width(50)

type action int
type itemKey string
type viewType int
type itemType int

const (
	NO_ACTION action = iota
	CREATE_TABLE
)

const (
	SINGLE_CHOICE viewType = iota
	MULTI_CHOICE
	NAVIGATION_CHOICE
	INPUT
	END
)

const (
	STANDARD itemType = iota
	HEADING
	CHOICE
)

type model struct {
	items        map[itemKey]*view // items on the to-do list
	cursor       int               // which to-do list item our cursor is pointing at
	selectedItem itemKey

	input      textinput.Model
	selection  string
	selections map[int]string

	err error

	// // below is state that covers all the posible Permutations of the entire app

	// // add / modify table
	// tableName  textinput.Model
	// isNewTable string
	// fields     []fieldSetting

	// // remove a table
	// removeTableNames []string

	// // generate things
	// tableNames []string
	// actions    []string // actions to execute at end, this gives us an option to build a series of actions.
}

type item struct {
	text     string
	itemType itemType
	key      string
	action   action
	question string
	style    lipgloss.Style
}

func (m *model) currentItem() *view {
	// logrus.Info("m.selectedItem", m.selectedItem)
	return m.items[m.selectedItem]
}

func (m *model) changeView(itemKey itemKey, action action) (tea.Model, tea.Cmd) {
	oldItem := m.currentItem()
	if oldItem != nil && oldItem.process != nil {
		oldItem.process(m)
	}
	m.selectedItem = itemKey

	if action != NO_ACTION {
		localstate.actions = append(localstate.actions, action)
	}
	m.cursor = 1
	newItem := m.currentItem()
	if newItem.viewType == INPUT {
		// current.input.
		m.input.SetValue("")
		m.input.SetCursorMode(textinput.CursorBlink)
		m.input.Focus()
		return m, textinput.Blink
	} else if newItem.viewType == MULTI_CHOICE || newItem.viewType == NAVIGATION_CHOICE || newItem.viewType == SINGLE_CHOICE {
		m.currentItem().items = m.currentItem().loadItems()
	} else if newItem.viewType == END {
		err := run()
		if err != nil {
			m.err = err
			return m, nil
		}
		return m, tea.Quit
	}
	return m, nil
}

/// TODO break this down in to composeable peices based on item type
type view struct {
	items     []*item
	loadItems func() []*item
	viewType  viewType
	viewKey   itemKey
	ref       interface{}
	question  string

	// changable pieces
	// selections map[int]string
	// selection  string
	// inputModel textinput.Model
	process func(*model)
}

type state struct {
	actions    []action
	tableName  string
	fieldNames []string
	fieldTypes []string
}

func (s *state) fields() []Field {
	fields := make([]Field, 0)
	for i, _ := range localstate.fieldNames {
		field := Field{
			FieldName:    localstate.fieldNames[i],
			FieldType:    localstate.fieldTypes[i],
			FieldDefault: "",
		}
		fields = append(fields, field)
	}
	return fields
}

var localstate *state

func init() {
	localstate = &state{
		actions:    make([]action, 0),
		tableName:  "",
		fieldNames: make([]string, 0),
		fieldTypes: make([]string, 0),
	}
}

// have different view types that resolve the action differently based on a view type e.g. choices are an input and have a default next key

// push state on to a slice that you can restore from at any point

// have a validate message

// have an exit

func initialModel() *model {
	model := &model{
		items: map[itemKey]*view{
			"home": {
				loadItems: func() []*item {
					return []*item{
						{text: "Choices choices", itemType: HEADING},
						{text: "Make database changes", key: "database"},
						{text: "Generate REST & UI templates"},
					}
					//Enter key : "enterTableName"
				},
				viewType: NAVIGATION_CHOICE,
			},
			"alt": {
				loadItems: func() []*item {
					return []*item{
						{text: "Database", itemType: HEADING},
						{text: "New table"},
						{text: "Add field to existing table"},
						{text: "> Rest"},
						{text: "Model"},
						{text: "Actions"},
						{text: "> User Interface"},
						{text: "List Page"},
						{text: "Edit Page"},
						{text: "TS Definition"},
						{text: "API"},
					}
					//Enter key : "enterTableName"
				},
				viewType: NAVIGATION_CHOICE,
			},
			"database": {
				loadItems: func() []*item {
					return []*item{
						{text: "Database", itemType: HEADING},
						{text: "New Table", key: "enterTableName", action: CREATE_TABLE},
						{text: "Delete Table", key: "selectTables"},
					}
				},
				viewType: NAVIGATION_CHOICE,
			},
			"selectTables": {
				loadItems: func() []*item {
					items := make([]*item, 0)
					rows := make([]string, 0)
					getDBConnection().DB.Select(&rows, `
						select table_name from information_schema.tables
						where table_schema = 'public'
					`)

					for _, row := range rows {
						items = append(items, &item{
							text:     row,
							key:      row,
							itemType: CHOICE,
						})
					}
					return items
				},
				viewType: MULTI_CHOICE,
			},
			"enterTableName": {
				viewType: INPUT,
				question: "Table Name?",
				viewKey:  "enterFieldName",
				process: func(m *model) {
					localstate.tableName = m.input.Value()
				},
			},
			"enterFieldName": {
				viewType: INPUT,
				question: "Field Name?",
				viewKey:  "selectFieldType",
				process: func(m *model) {
					localstate.fieldNames = append(localstate.fieldNames, m.input.Value())
				},
			},
			"selectFieldType": {
				loadItems: func() []*item {
					return []*item{
						{text: "ULID", key: "character varying(26)"},
						{text: "Text", key: "text"},
						{text: "Numeric", key: "numeric"},
						{text: "DateTime", key: "timestamptz"},
						{text: "FK", key: "foreign_key"},
					}
				},
				// selection: "",
				viewType: SINGLE_CHOICE,
				viewKey:  "loopFIeld",
				process: func(m *model) {
					localstate.fieldTypes = append(localstate.fieldTypes, m.selection)
				},
			},
			"loopFIeld": {
				loadItems: func() []*item {
					return []*item{
						{text: "Keep going?", itemType: HEADING},
						{text: "Another Field", key: "enterFieldName"},
						{text: "Done", key: "end"},
					}
				},
				viewType: NAVIGATION_CHOICE,
			},
			"end": {
				viewType: END,
			},
		},

		cursor: 1,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selectedItem: "",
		selection:    "",
		selections:   make(map[int]string),
		input:        textinput.NewModel(),
	}

	model.changeView("home", NO_ACTION)
	return model
}

func main() {
	// load from .env file where scaffold is run
	if err := dotenv.Load(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		if m.currentItem().viewType == INPUT {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyEsc {
				return m.changeView(m.currentItem().viewKey, NO_ACTION)
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		choices := m.currentItem().items
		if choices != nil && len(choices) > 0 {
			choice := choices[m.cursor]

			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--

					if choices[m.cursor].itemType == HEADING {
						if m.cursor > 0 {
							m.cursor--
						} else {
							m.cursor++ // put cursor back
						}
					}
				}

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.cursor < len(choices)-1 {
					m.cursor++
					if m.cursor < len(choices)-1 && choices[m.cursor].itemType == HEADING {
						m.cursor++
					}
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case " ", "enter":
				if msg.String() == " " && m.currentItem().viewType == MULTI_CHOICE {
					selections := m.selections
					_, ok := selections[m.cursor]
					if ok {
						delete(selections, m.cursor)
					} else {
						selections[m.cursor] = choice.key
					}
				} else if m.currentItem().viewType == SINGLE_CHOICE {
					m.selection = choice.key
					return m.changeView(m.currentItem().viewKey, choice.action)
				} else if m.currentItem().viewType == NAVIGATION_CHOICE {
					return m.changeView(itemKey(choice.key), choice.action)
				}
			}
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *model) View() string {
	viewType := m.currentItem().viewType
	s := ""
	if viewType == END {
		if m.err != nil {
			s += fmt.Sprintf(
				"%s\n\n%s",
				m.err.Error(),
				"(esc to quit)",
			) + "\n"
		}
	} else if viewType == INPUT {
		s += fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			m.currentItem().question,
			m.input.View(),
			"(esc to quit)",
		) + "\n"
	} else if viewType == SINGLE_CHOICE || viewType == MULTI_CHOICE || viewType == NAVIGATION_CHOICE {
		choices := m.currentItem().items

		// Iterate over our choices
		for i, choice := range choices {

			// Is the cursor pointing at this choice?
			isHighlighted := m.cursor == i

			// Is this choice selected?
			checked := " " // not selected
			if m.selections != nil {
				selections := m.selections
				if _, ok := selections[i]; ok {
					checked = "+" // selected!
				}
			}

			line := choice.text
			if choice.itemType == HEADING {
				line = headingStyle.Render(choice.text) + BREAK
			} else if choice.itemType == CHOICE {
				line = fmt.Sprintf("[%s] %s", checked, choice.text)
			}

			if isHighlighted {
				s += selectionStyle.Render(line) + BREAK
			} else {
				s += line + BREAK
			}

			// Render the row
		}
	}

	// Send the UI for rendering
	return s
}

func run() (err error) {
	render := getRenderer()
	for _, action := range localstate.actions {
		if action == CREATE_TABLE {
			err = createTable(localstate.tableName, localstate.fields(), render, getDBConnection())
		}
	}
	return err
}

// // needs refactor
// func runActions(actions []string) {
// 	app := cli.NewApp()
// 	app.Name = "builder"
// 	app.Usage = "code generation"
// 	flags := []cli.Flag{
// 		cli.StringFlag{
// 			Name:  "skip",
// 			Usage: "skip bcomp diff",
// 		},
// 	}
// 	// app.Flags = flags
// 	app.Commands = []cli.Command{
// 		{
// 			Name:    "scaffold",
// 			Aliases: []string{"s"},
// 			Usage:   "Create a new project [projectname] [project path]",
// 			Action: func(c *cli.Context) error {
// 				return createProject(c, render)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "table",
// 			Aliases: []string{"t"},
// 			Usage:   "Create a new table [tablename] [fieldname:fieldtype:fielddefault]",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					return err
// 				}
// 				db := getDBConnection()
// 				return createTable(c, render, db)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "fields",
// 			Aliases: []string{"f"},
// 			Usage:   "Add fields to an existing table [tablename] [fieldname:fieldtype:fielddefault]",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					return err
// 				}
// 				db := getDBConnection()
// 				return addFields(c, render, db)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "model",
// 			Aliases: []string{"m"},
// 			Usage:   "create a model from a table [tablename]",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					return err
// 				}
// 				db := getDBConnection()
// 				return createModel(c, render, db)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "rest",
// 			Aliases: []string{"r"},
// 			Usage:   "create a restful interface from a table [tablename]",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					return err
// 				}
// 				db := getDBConnection()
// 				return createRest(c, render, db)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "edit",
// 			Aliases: []string{"e"},
// 			Usage:   "create a edit page from a table [tablename]",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					return err
// 				}
// 				db := getDBConnection()
// 				return createEdit(c, render, db)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "list",
// 			Aliases: []string{"l"},
// 			Usage:   "create a list page from a table [tablename]",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					return err
// 				}
// 				db := getDBConnection()
// 				return createList(c, render, db)
// 			},
// 			Flags: flags,
// 		},
// 		{
// 			Name:    "migration",
// 			Aliases: []string{"mi"},
// 			Usage:   "perform schema migration",
// 			Action: func(c *cli.Context) error {
// 				err := dotenv.Load() // load from .env file where scaffold is run
// 				if err != nil {
// 					panic(err)
// 				}
// 				db := getDBConnection()
// 				err = doMigration(c, render, db)
// 				if err != nil {
// 					panic(err)
// 				}
// 				return err
// 			},
// 			Flags: flags,
// 		},
// 	}
// 	app.Run(os.Args)
// }
