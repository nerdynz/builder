package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dotenv "github.com/joho/godotenv"
	errors "github.com/kataras/go-errors"
	"github.com/sirupsen/logrus"
)

const (
	maxWidth = 40
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

var subheadingStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#282D3F")).
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
type viewItemKey string
type viewType string
type itemType int

const (
	NO_ACTION action = iota
	SCAFFOLD_PROJECT
	CREATE_TABLE
	ADD_FIELDS
	GENERATE_MODEL
	GENERATE_REST
	GENERATE_EDIT
	GENERATE_LIST
	GENERATE_SEARCH
	GENERATE_MIGRATION
	MIGRATE
)

const (
	SINGLE_CHOICE     viewType = "SINGLE_CHOICE"
	MULTI_CHOICE      viewType = "MULTI_CHOICE"
	NAVIGATION_CHOICE viewType = "NAVIGATION_CHOICE"
	PRIORITY_CHOICE   viewType = "PRIORITY_CHOICE"
	INPUT             viewType = "INPUT"
	PROGRESS          viewType = "PROGRESS"
	END               viewType = "END"
)

const (
	STANDARD itemType = iota
	HEADING
	CHOICE
	PRIORITY
)

/// TODO break this down in to composeable peices based on item type
type view struct {
	items     []*item
	loadItems func() []*item
	viewType  viewType
	viewKey   viewItemKey
	ref       interface{}
	question  string

	// changable pieces
	// selections map[int]string
	// selection  string
	// inputModel textinput.Model
	process func(*model)
}

type state struct {
	actions []action
	// PROJECT
	projectName string
	projectPath string

	// DATABASE
	tableName       string
	tables          []string
	fieldNames      []string
	fieldTypes      []string
	fieldPriorities []string
}

func (s *state) fields() []Field {
	fields := make([]Field, 0)
	for i := range localstate.fieldNames {
		fieldType := ""
		fieldPriority := ""
		if len(localstate.fieldTypes) > i {
			fieldType = localstate.fieldTypes[i]
		}
		if len(localstate.fieldPriorities) > i {
			fieldPriority = localstate.fieldPriorities[i]
		}

		field := Field{
			FieldName:     localstate.fieldNames[i],
			FieldType:     fieldType,
			FieldPriority: fieldPriority,
			FieldDefault:  "",
		}
		fields = append(fields, field)
	}
	return fields
}

type model struct {
	items        map[viewItemKey]*view // items on the to-do list
	cursor       int                   // which to-do list item our cursor is pointing at
	selectedItem viewItemKey

	input              textinput.Model
	progress           progress.Model
	selection          string
	selections         map[int]string
	prioritySelections map[int]string

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

type tickMsg time.Time
type errMsg error

var localstate *state

func (m *model) currentItem() *view {
	// logrus.Info("m.selectedItem", m.selectedItem)
	return m.items[m.selectedItem]
}

func (m *model) changeView(itemKey viewItemKey, action action) (tea.Model, tea.Cmd) {
	_, ok := m.items[itemKey]
	if !ok {
		m.err = errors.New("item doesn't exist " + string(itemKey) + " doesnt exist")
		m.changeView("error", NO_ACTION)
	}
	m.progress.SetPercent(0)
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
	if newItem.viewType == END {
		// maybe later
	} else if newItem.viewType == INPUT {
		// current.input.
		m.input.SetValue("")
		m.input.SetCursorMode(textinput.CursorBlink)
		m.input.Focus()
		return m, textinput.Blink
	} else if newItem.viewType == MULTI_CHOICE || newItem.viewType == NAVIGATION_CHOICE || newItem.viewType == SINGLE_CHOICE || newItem.viewType == PRIORITY_CHOICE {
		m.currentItem().items = m.currentItem().loadItems()
	} else if newItem.viewType == PROGRESS {
		go func() {
			err := run()
			if err != nil {
				m.err = err
				m.changeView("error", NO_ACTION)
			}
		}()
		return m, nil
	}
	return m, nil
}

// have different view types that resolve the action differently based on a view type e.g. choices are an input and have a default next key

// push state on to a slice that you can restore from at any point

// have a validate message

// have an exit

func initialModel() *model {
	model := &model{
		items: map[viewItemKey]*view{
			"home": {
				loadItems: func() []*item {
					items := []*item{
						{text: "Choices choices", itemType: HEADING},
						{text: "Create New Project", key: "enterProjectName", action: SCAFFOLD_PROJECT},
					}

					if isEnvPresent {
						items = append(items, &item{text: "Make database changes", key: "database"})
						items = append(items, &item{text: "Generate REST & UI templates", key: "alt"})
						items = append(items, &item{text: "Create Search", key: "selectTableForSearch"})
					}

					return items
					//Enter key : "enterTableName"
				},
				viewType: NAVIGATION_CHOICE,
			},
			"error": {
				viewType: END,
			},
			"alt": {
				loadItems: func() []*item {
					return []*item{
						{text: "Model", action: GENERATE_MODEL, key: "selectTables"},
						{text: "Actions", action: GENERATE_REST, key: "selectTables"},
						{text: "List Page", action: GENERATE_LIST, key: "selectTables"},
						{text: "Edit Page", action: GENERATE_EDIT, key: "selectTables"},
						// {text: "TS Definition", action: GENERATE_EDIT
						// {text: "API"},
					}
					//Enter key : "enterTableName"
				},
				viewType: NAVIGATION_CHOICE,
			},
			"enterProjectName": {
				viewType: INPUT,
				question: "Project Name",
				viewKey:  "enterProjectPath",
				process: func(m *model) {
					localstate.projectName = m.input.Value()
				},
			},
			"enterProjectPath": {
				viewType: INPUT,
				question: "Project Path",
				viewKey:  "end",
				process: func(m *model) {
					localstate.projectPath = m.input.Value()
				},
			},
			"database": {
				loadItems: func() []*item {
					return []*item{
						{text: "Database", itemType: HEADING},
						{text: "New Table", key: "enterTableName", action: CREATE_TABLE},
						{text: "Delete Table", key: "selectTables"},
						{text: "Migrate", key: "end", action: MIGRATE},
					}
				},
				viewType: NAVIGATION_CHOICE,
			},
			"selectTable": {
				process: func(m *model) {
					localstate.tableName = m.selection
				},
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
				viewKey:  "end",
				viewType: SINGLE_CHOICE,
			},
			"selectTableForSearch": {
				process: func(m *model) {
					localstate.tableName = m.selection
				},
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
							itemType: STANDARD,
							action:   GENERATE_SEARCH,
						})
					}
					return items
				},
				viewKey:  "assignSearchPriority",
				viewType: SINGLE_CHOICE,
			},
			"assignSearchPriority": {
				process: func(m *model) {
					localstate.fieldNames = make([]string, 0)
					localstate.fieldPriorities = make([]string, 0)
					for i, fieldName := range m.selections {
						localstate.fieldNames = append(localstate.fieldNames, fieldName)
						localstate.fieldPriorities = append(localstate.fieldPriorities, m.prioritySelections[i])
					}
				},
				loadItems: func() []*item {
					items := make([]*item, 0)
					rows := make([]string, 0)
					getDBConnection().DB.Select(&rows, `
					select column_name from information_schema.columns
					where table_schema = 'public'
					and table_name = '`+localstate.tableName+`'
					and column_name <> 'tsv'
					order by column_name
					`)

					for _, row := range rows {
						items = append(items, &item{
							text:     row,
							key:      row,
							itemType: PRIORITY,
						})
					}
					return items
				},
				viewType: PRIORITY_CHOICE,
				viewKey:  "end",
			},
			"selectTables": {
				process: func(m *model) {
					tables := make([]string, 0)
					for _, selection := range m.selections {
						tables = append(tables, selection)
					}
					localstate.tables = tables
				},
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
				viewKey:  "end",
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
				viewKey:  "loopField",
				process: func(m *model) {
					localstate.fieldTypes = append(localstate.fieldTypes, m.selection)
				},
			},
			"loopField": {
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
				viewType: PROGRESS,
			},
		},

		cursor: 1,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selectedItem:       "",
		selection:          "",
		selections:         make(map[int]string),
		prioritySelections: make(map[int]string),
		input:              textinput.NewModel(),
		progress:           progress.NewModel(progress.WithDefaultGradient()),
	}

	model.changeView("home", NO_ACTION)
	return model
}

var isEnvPresent bool

func main() {
	isEnvPresent = false
	// load from .env file where scaffold is run
	if err := dotenv.Load(); err == nil {
		isEnvPresent = true
	}
	if err := dotenv.Load("./rest/.env"); err == nil {
		isEnvPresent = true
	}

	// if !isEnvPresent {
	// 	fmt.Printf("Failed to load .env file")
	// 	os.Exit(1)
	// }

	localstate = &state{
		actions:    make([]action, 0),
		tableName:  "",
		fieldNames: make([]string, 0),
		fieldTypes: make([]string, 0),
	}

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *model) Init() tea.Cmd {
	return tickCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil
	// case tea.WindowSizeMsg:
	// 	m.progress.Width = msg.Width - padding*2 - 4
	// 	if m.progress.Width > maxWidth {
	// 		m.progress.Width = maxWidth
	// 	}
	// 	return m, nil
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case tickMsg:
		if m.progress.Percent() == 1.00 {
			return m, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.0)
		if m.currentItem().viewType == PROGRESS {
			cmd = m.progress.IncrPercent(0.10)
		}

		return m, tea.Batch(tickCmd(), cmd)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		if m.currentItem().viewType == INPUT {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyEsc {
				return m.changeView(m.currentItem().viewKey, NO_ACTION)
			}
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

			case "a", "b", "c", "d", "e", "f", "backspace":
				if m.currentItem().viewType == PRIORITY_CHOICE {
					_, ok := m.selections[m.cursor]
					if ok && (m.prioritySelections[m.cursor] == msg.String() || msg.String() == "backspace") { // same key as already there or delete button
						delete(m.selections, m.cursor)
						delete(m.prioritySelections, m.cursor)
					} else {
						m.selections[m.cursor] = choice.key
						m.prioritySelections[m.cursor] = msg.String()
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
					return m.changeView(viewItemKey(choice.key), choice.action)
				} else {
					return m.changeView(m.currentItem().viewKey, choice.action)
				}
			default:
				logrus.Info(msg.String())
			}
		}
	default:
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	if m.currentItem().viewType == INPUT {
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *model) View() string {
	if m.err != nil {
		return headingStyle.Render("ERROR") + BREAK + m.err.Error()
	}

	viewType := m.currentItem().viewType
	s := string(viewType) + BREAK
	// return "\n" +
	// pad + helpStyle("Press any key to quit")
	if viewType == END {
		if m.err != nil {
			s += fmt.Sprintf(
				"%s\n\n%s",
				m.err.Error(),
				"(esc to quit)",
			) + "\n"
		}
	} else if viewType == PROGRESS {
		s += m.progress.View() + BREAK
	} else if viewType == INPUT {
		s += fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			m.currentItem().question,
			m.input.View(),
			"(esc to quit)",
		) + "\n"
	} else if viewType == SINGLE_CHOICE || viewType == MULTI_CHOICE || viewType == NAVIGATION_CHOICE || viewType == PRIORITY_CHOICE {
		choices := m.currentItem().items

		// if viewType == PRIORITY_CHOICE {
		// 	logrus.Info(choices)
		// }
		// Iterate over our choices
		for i, choice := range choices {

			// Is the cursor pointing at this choice?
			isHighlighted := m.cursor == i

			line := choice.text
			if choice.itemType == HEADING {
				line = subheadingStyle.Render(choice.text) + BREAK // IGNORING
			} else if choice.itemType == CHOICE {
				// Is this choice selected?
				checked := " " // not selected
				if m.selections != nil {
					selections := m.selections
					if _, ok := selections[i]; ok {
						checked = "+" // selected!
					}
				}
				line = fmt.Sprintf("[%s] %s", checked, choice.text)
			} else if choice.itemType == PRIORITY {
				char := " "
				if m.selections != nil {
					if _, ok := m.selections[i]; ok {
						char = m.prioritySelections[i]
					}
				}
				line = fmt.Sprintf("[%s] %s", char, choice.text)
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
		if action == SCAFFOLD_PROJECT {
			err = createProject(localstate.projectName, localstate.projectPath)
		} else if action == CREATE_TABLE {
			err = createTable(localstate.tableName, localstate.fields(), render, getDBConnection())
		} else if action == ADD_FIELDS {
			for _, tableName := range localstate.tables {
				err = addFields(tableName, localstate.fields(), render, getDBConnection())
			}
		} else if action == GENERATE_MODEL {
			for _, tableName := range localstate.tables {
				err = createModel(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_REST {
			for _, tableName := range localstate.tables {
				err = createRest(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_EDIT {
			for _, tableName := range localstate.tables {
				err = createEdit(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_LIST {
			for _, tableName := range localstate.tables {
				err = createList(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_SEARCH {
			err = createSearch(localstate.tableName, localstate.fields(), render, getDBConnection())
		} else if action == MIGRATE {
			err = doMigration(render, getDBConnection())
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
