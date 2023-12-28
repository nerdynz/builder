package main

import (
	"fmt"
	"log"
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

var endOf = lipgloss.NewStyle().
	Background(lipgloss.Color("#15b8a6")).
	MarginTop(1).
	Width(50)

type action int
type viewItemKey string
type viewType string
type itemType int

const (
	NO_ACTION action = iota
	SCAFFOLD_PROJECT
	CREATE_TABLE
	DELETE_TABLE
	CREATE_MIGRATION
	ADD_FIELDS
	GENERATE_API
	GENERATE_MODEL
	GENERATE_PROTO
	GENERATE_RPC
	GENERATE_REST
	GENERATE_EDIT
	GENERATE_LIST
	GENERATE_BE
	GENERATE_EVERYTHING
	GENERATE_FE
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

// / TODO break this down in to composeable peices based on item type
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
	initialInputValue string
	tables            []string
	fieldNames        []string
	fieldTypes        []string
	fieldPriorities   []string
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
		// go func() {
		err := run()
		if err != nil {
			logrus.Error("err", err)
			m.err = err
			m.changeView("error", NO_ACTION)
		}
		// }()
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
						{text: "Create", itemType: HEADING},
					}

					if isEnvPresent {
						items = append(items, &item{text: "Database", key: "database"})
						items = append(items, &item{text: "Scaffold", key: "alt"})
					}
					if isEnvPresent {
						items = append(items, &item{text: "Migrate", itemType: HEADING})
						items = append(items, &item{text: "Run Migration", key: "end", action: MIGRATE})
					}

					items = append(items, &item{text: "New", itemType: HEADING})
					items = append(items, &item{text: "New Project", key: "enterProjectName", action: SCAFFOLD_PROJECT})
					return items
				},
				viewType: NAVIGATION_CHOICE,
			},
			"error": {
				viewType: END,
			},
			"alt": {
				loadItems: func() []*item {
					return []*item{

						{text: "Backend", itemType: HEADING},
						{text: "Proto", action: GENERATE_PROTO, key: "selectTables"},
						{text: "Model", action: GENERATE_MODEL, key: "selectTables"},
						{text: "RPC", action: GENERATE_RPC, key: "selectProtos"},
						// {text: "Actions", action: GENERATE_REST, key: "selectTables"},
						{text: "ALL", action: GENERATE_BE, key: "selectTables"},

						{text: "Frontend", itemType: HEADING},
						{text: "List Page", action: GENERATE_LIST, key: "selectTables"},
						{text: "Edit Page", action: GENERATE_EDIT, key: "selectTables"},
						{text: "API", action: GENERATE_API, key: "selectTables"},
						{text: "ALL", action: GENERATE_FE, key: "selectTables"},
						{text: "EVERYTHING", action: GENERATE_EVERYTHING, key: "selectTables"},
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
						{text: "Blank Migration", key: "enterMigrationDescription", action: CREATE_MIGRATION},
						{text: "New Table", key: "enterTableName", action: CREATE_TABLE},
						{text: "Delete Table", key: "selectTables", action: DELETE_TABLE},
						{text: "Create Search", key: "selectTableForSearch"},
					}
				},
				viewType: NAVIGATION_CHOICE,
			},
			"selectTable": {
				process: func(m *model) {
					localstate.initialInputValue = m.selection
				},
				loadItems: func() []*item {
					items := make([]*item, 0)
					rows := make([]string, 0)
					getDBConnection().DB.Select(&rows, `
						select table_name from information_schema.tables
						where table_schema = 'public'
					`)
					items = append(items, &item{text: "Select a Table", itemType: HEADING})

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
					localstate.initialInputValue = m.selection
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
					and table_name = '`+localstate.initialInputValue+`'
					and column_name <> 'tsv'
					order by column_name
					`)
					items = append(items, &item{text: "Assign Priority to Search", itemType: HEADING})

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
					items = append(items, &item{text: "Select Tables", itemType: HEADING})

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
			"selectProtos": {
				process: func(m *model) {
					tables := make([]string, 0)
					for _, selection := range m.selections {
						tables = append(tables, selection)
					}
					localstate.tables = tables
				},
				loadItems: func() []*item {
					items := make([]*item, 0)
					files, err := os.ReadDir("./proto")
					if err != nil {
						log.Fatal(err)
					}

					// for _, f := range files {
					// 	fmt.Println(f.Name())
					// }
					items = append(items, &item{text: "Select Protos", itemType: HEADING})
					for _, f := range files {
						items = append(items, &item{
							text:     f.Name(),
							key:      f.Name(),
							itemType: CHOICE,
						})
					}
					return items
				},
				viewKey:  "end",
				viewType: MULTI_CHOICE,
			},
			"enterMigrationDescription": {
				viewType: INPUT,
				question: "Description of Migration",
				viewKey:  "end",
				process: func(m *model) {
					localstate.initialInputValue = m.input.Value()
				},
			},
			"enterTableName": {
				viewType: INPUT,
				question: "Table Name?",
				viewKey:  "enterFieldName",
				process: func(m *model) {
					localstate.initialInputValue = m.input.Value()
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
	if err := dotenv.Load(".builder.env"); err == nil {
		isEnvPresent = true
	}
	for i := 1; i <= 5; i++ {
		if isEnvPresent {
			break
		} else {
			if err := os.Chdir(".."); err != nil {
				logrus.Fatal("Chdir failed", err)
			}
			pw, _ := os.Getwd()
			logrus.Info("WD =>", pw)

			if err := dotenv.Load(".builder.env"); err == nil {
				isEnvPresent = true
			}
		}
	}

	// if err := dotenv.Load("./rpc/.env"); err == nil {
	// 	isEnvPresent = true
	// }

	// dotenv.Load(".builder.env")
	// if !isEnvPresent {
	// 	logrus.Fatal("Failed to load .env file")
	// }

	localstate = &state{
		actions:           make([]action, 0),
		initialInputValue: "",
		fieldNames:        make([]string, 0),
		fieldTypes:        make([]string, 0),
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
	s := ""
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
				line = headingStyle.Render(choice.text) + BREAK // IGNORING
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

	s += endOf.Render("") // IGNORING
	// Send the UI for rendering
	return s
}

func run() (err error) {
	render := getRenderer()
	for _, action := range localstate.actions {
		if action == SCAFFOLD_PROJECT {
			err = createProject(localstate.projectName, localstate.projectPath)
		} else if action == CREATE_MIGRATION {
			err = createBlankMigration(localstate.initialInputValue, render, getDBConnection())
		} else if action == CREATE_TABLE {
			err = createTable(localstate.initialInputValue, localstate.fields(), render, getDBConnection())
		} else if action == ADD_FIELDS {
			for _, tableName := range localstate.tables {
				err = addFields(tableName, localstate.fields(), render, getDBConnection())
			}
		} else if action == GENERATE_API {
			for _, tableName := range localstate.tables {
				err = createAPI(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_MODEL {
			for _, tableName := range localstate.tables {
				err = createModel(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_PROTO {
			for _, tableName := range localstate.tables {
				err = createProto(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_RPC {
			for _, tableName := range localstate.tables {
				err = createRPC(tableName, render, getDBConnection())
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
		} else if action == GENERATE_BE {
			for _, tableName := range localstate.tables {
				if err = createProto(tableName, render, getDBConnection()); err != nil {
					return err
				}
				if err = createModel(tableName, render, getDBConnection()); err != nil {
					return err
				}
				if err = createRPC(tableName, render, getDBConnection()); err != nil {
					return err
				}
			}
		} else if action == GENERATE_FE {
			for _, tableName := range localstate.tables {
				err = createAPI(tableName, render, getDBConnection())
				err = createEdit(tableName, render, getDBConnection())
				err = createList(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_BE {
			for _, tableName := range localstate.tables {
				if err = createProto(tableName, render, getDBConnection()); err != nil {
					return err
				}
				if err = createModel(tableName, render, getDBConnection()); err != nil {
					return err
				}
				if err = createRPC(tableName, render, getDBConnection()); err != nil {
					return err
				}
			}
			for _, tableName := range localstate.tables {
				err = createAPI(tableName, render, getDBConnection())
				err = createEdit(tableName, render, getDBConnection())
				err = createList(tableName, render, getDBConnection())
			}
		} else if action == GENERATE_SEARCH {
			err = createSearch(localstate.initialInputValue, localstate.fields(), render, getDBConnection())
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
