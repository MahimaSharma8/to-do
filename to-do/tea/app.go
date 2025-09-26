package tea 

import (
	"fmt"
	"strings"
	"time"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/MahimaSharma8/to-do/todo"
)


type step int 
const (
    stepList step = iota
    stepTask
    stepTopic
    stepStatus
    stepPriority
    stepDone
)

const todoFile = "./.todos.json"
var appBox = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#FF69B4")). // pink border
    Padding(1, 2).
    Margin(1, 2)


//for tab customization
var (
    highlightColor   = lipgloss.AdaptiveColor{Light: "#f4e1ea", Dark: "#FF69B4"} // pink
    inactiveTabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Foreground(lipgloss.Color("#fdc9e3ff")).
        Bold(true)
    activeTabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Foreground(lipgloss.Color("230")).
        Background(lipgloss.Color("#f41283ff")). // hot pink highlight
        Bold(true)
)



var todoBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(60)

var (
    itemStyle   = lipgloss.NewStyle().Padding(0, 2)
    selectedItem= lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true).Underline(true)
)



func renderTabs(m model) string {
    tabs := []string{"List", "Task", "Topic", "Status", "Priority", "Done"}
    var renderedTabs []string
    for i, t := range tabs {
        if i == int(m.step) {
            renderedTabs = append(renderedTabs, activeTabStyle.Render(" "+t+" "))
        } else {
            renderedTabs = append(renderedTabs, inactiveTabStyle.Render(" "+t+" "))
        }
    }
    return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}


func renderTodos(todos *todo.Todos, selected int) string {
	if len(*todos) == 0 {
		return lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("238")).
			Render("No todos yet!")
	}

	var rendered []string
	for i, t := range *todos {
		// Prefix based on status
		prefix := "[ ]"
		if t.Status == todo.Done {
			prefix = "[x]"
		}

		// Task line
		taskLine := lipgloss.NewStyle().
			Bold(true).
			Render(fmt.Sprintf("%s %s", prefix, t.Task))

		details := []string{
			fmt.Sprintf("Topic   : %s", t.Topic),
			fmt.Sprintf("Status  : %s", t.Status),
			fmt.Sprintf("Priority: %d", t.Priority),
		}
		if t.Completed != nil {
			details = append(details, fmt.Sprintf("Time Spent: %s", t.TimeWorked.Truncate(time.Second)))
		}

		content := lipgloss.JoinVertical(lipgloss.Left, append([]string{taskLine}, details...)...)

		boxStyle := todoBox
		if i == selected {
			boxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FF69B4")).
				Padding(0, 1)
		}

		rendered = append(rendered, boxStyle.Render(content))
	}

	return strings.Join(rendered, "\n\n")
}




func renderChoiceList(choices []string, selected int) string {
    var rendered []string
    for i, choice := range choices {
        if i == selected {
            rendered = append(rendered, selectedItem.Render("> "+choice))
        } else {
            rendered = append(rendered, itemStyle.Render("  "+choice))
        }
    }

    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#FF69B4")).
        Padding(1, 2).
        Render(strings.Join(rendered, "\n"))
}



type model struct {
	step step
	totalSteps int
	taskInput textinput.Model

	selectedTopic    int     //so these are like indexes to the string array
	selectedStatus   int
	selectedPriority int

	topics     []string
	statuses   []string
	priorities []string

	tabs      []string

	editing bool

	selectedTodo int

	todos *todo.Todos
}

func parsePriority(s string) todo.Priority {
	switch s {
	case "Low":
		return todo.Low
	case "Medium":
		return todo.Medium
	case "High":
		return todo.High
	default:
		return todo.Low
	}
}

func parseStatus(s string) todo.Status {
	switch s {
	case "Pending":
		return todo.Pending
	case "In-progress":
		return todo.InProgress
	case "Done":
		return todo.Done
	default:
		return todo.Pending
	}
}




func initialModel() model {
    // Load todos from file
    var todos todo.Todos
    _ = todos.Load(todoFile)

    // Initialize task input
    ti := textinput.New()
    ti.Placeholder = "What are we doing today?"
    ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true)
    ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
    ti.Focus() // start focused

    // Initialize progress bar with gradient
    p := progress.New(
        progress.WithGradient("#FF69B4", "#FF1493"),
        progress.WithDefaultGradient(),
    )
    p.Width = 50

    // Construct model
    m := model{
        step:       stepList,
        totalSteps: 5,
        taskInput:  ti,
        topics:     []string{"Work", "Personal", "Health", "Study"},
        statuses:   []string{"Pending", "In-progress", "Done"},
        priorities: []string{"Low", "Medium", "High"},
        tabs:       []string{"Add", "List"},
        todos:      &todos,
    }

    return m
}



func (m model) Init() tea.Cmd {   //Init can return a Cmd that could perform some initial I/O.
	return tea.Batch(
		textinput.Blink,
	) //blink for input

}


//The update function is called when ”things happen.” Its job is to look at what has happened and return an updated model in response. It can also return a Cmd to make more things happen

//The “something happened” comes in the form of a Msg, which can be any type. Messages are the result of some I/O that took place, such as a keypress, timer tick, or a response from a server.


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.step == stepTask {
			var cmd tea.Cmd
			m.taskInput, cmd = m.taskInput.Update(msg)
			// Move to next step if Enter is pressed
			if msg.String() == "enter" {
				m.step++
			}
			return m, cmd
		}
		switch msg.String() {
		// Quit keys
		case "ctrl+c", "q":
			_ = m.todos.Store(todoFile)
			return m, tea.Quit
		case "up":
			switch m.step {
			case stepList:
				if m.selectedTodo > 0 {
					m.selectedTodo--
				}
			case stepTopic:
				if m.selectedTopic > 0 {
					m.selectedTopic--
				}
			case stepStatus:
				if m.selectedStatus > 0 {
					m.selectedStatus--
				}
			case stepPriority:
				if m.selectedPriority > 0 {
					m.selectedPriority--
				}
			}
			return m, nil

    	case "down":
			switch m.step {
			case stepList:
				if m.selectedTodo < len(*m.todos)-1 {
					m.selectedTodo++
				}
			case stepTopic:
				if m.selectedTopic < len(m.topics)-1 {
					m.selectedTopic++
				}
			case stepStatus:
				if m.selectedStatus < len(m.statuses)-1 {
					m.selectedStatus++
				}
			case stepPriority:
				if m.selectedPriority < len(m.priorities)-1 {
					m.selectedPriority++
				}
			}
			return m, nil
		case "right":
			if m.step < stepDone {
				m.step++
			}

		case "left":
			if m.step > stepList {
				m.step--
			}
		case "e":
			if m.step == stepList && len(*m.todos) > 0 {
				m.editing = true
				m.step = stepTask
				m.taskInput.SetValue((*m.todos)[m.selectedTodo].Task)
			}
			return m, nil
		case "x":
			if m.step == stepList && len(*m.todos) > 0 {
				// Mark selected todo as done
				err := m.todos.Complete(m.selectedTodo + 1) // your Complete method expects 1-based index
				if err == nil {
					_ = m.todos.Store(todoFile)
				}
				return m, nil
			}

		case "d":
			if m.step == stepList && len(*m.todos) > 0 {
				// Delete selected todo
				err := m.todos.Delete(m.selectedTodo + 1) // 1-based
				if err == nil {
					// Adjust selection if we were at the end
					if m.selectedTodo >= len(*m.todos) {
						m.selectedTodo = len(*m.todos) - 1
					}
					_ = m.todos.Store(todoFile)
				}
				return m, nil
			}
		case "enter":
		// Editing an existing todo
		if m.editing && m.step == stepTask && m.selectedTodo >= 0 {
			_ = m.todos.Edit(
				m.selectedTodo,
				m.taskInput.Value(),
				string((*m.todos)[m.selectedTodo].Status),
				int((*m.todos)[m.selectedTodo].Priority),
				(*m.todos)[m.selectedTodo].Topic,
			)
			_ = m.todos.Store(todoFile)
			m.step = stepList
			m.taskInput.SetValue("")
			m.editing = false
			return m, nil
		}

		// Creating a new todo via wizard
		if !m.editing {
			switch m.step {
			case stepTask:
				// Move to next step in wizard
				m.step = stepTopic
			case stepTopic:
				m.step = stepStatus
			case stepStatus:
				m.step = stepPriority
			case stepPriority:
				m.step = stepDone
			case stepDone:
				// Save todo at the end
				task := m.taskInput.Value()
				status := parseStatus(m.statuses[m.selectedStatus])
				priority := parsePriority(m.priorities[m.selectedPriority])
				topic := m.topics[m.selectedTopic]

				m.todos.Add(task, status, priority, topic)
				_ = m.todos.Store(todoFile)

				// Reset
				m.step = stepList
				m.taskInput.SetValue("")
			}
		}
			return m, nil
		
		}
	}

	// ---- Let text input handle typing when on stepTask ----
	if m.step == stepTask {
		var cmd tea.Cmd
		m.taskInput, cmd = m.taskInput.Update(msg)
		return m, cmd
	}

	return m, nil
}


//We look at the model in its current state and use it to return a string. That string is our UI!

func (m model) View() string {
    var body string

    switch m.step {
    case stepList:
        body = lipgloss.NewStyle().Bold(true).Render("Your Todos:") + "\n\n" + renderTodos(m.todos, m.selectedTodo)
    case stepTask:
        body = "What are we doing today?\n\n" + m.taskInput.View()
    case stepTopic:
        body = "Choose a topic:\n\n" + renderChoiceList(m.topics, m.selectedTopic)
    case stepStatus:
        body = "Pick a status:\n\n" + renderChoiceList(m.statuses, m.selectedStatus)
    case stepPriority:
        body = "Select priority:\n\n" + renderChoiceList(m.priorities, m.selectedPriority)
    case stepDone:
        body = "Okay, on it now!\n\n" +
            fmt.Sprintf("Task: %s\n", m.taskInput.Value()) +
            fmt.Sprintf("Topic: %s\n", m.topics[m.selectedTopic]) +
            fmt.Sprintf("Status: %s\n", m.statuses[m.selectedStatus]) +
            fmt.Sprintf("Priority: %s\n", m.priorities[m.selectedPriority]) +
            "\nPress Enter to save.\n"
    }

    help := lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Render("'<' and '>' to move \tx to mark done\te to edit \t d to delete\t q to quit")


    // Wrap *everything* inside the pink app box
    return appBox.Render(
        renderTabs(m) + "\n\n" + body + "\n\n" + help + "\n\n",
    )
}


// StartApp is the entrypoint from cmd/main.go
func StartApp() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}