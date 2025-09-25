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

func tabBorderWithBottom(left, mid, right string) lipgloss.Border {
    return lipgloss.Border{
        Top:         "─",
        Bottom:      mid,
        Left:        left,
        Right:       right,
        TopLeft:     "┌",
        TopRight:    "┐",
        BottomLeft:  left,
        BottomRight: right,
    }
}

//for tab customization
var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#FF69B4", Dark: "#FF69B4"} // pinkish!
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)

	tabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Foreground(lipgloss.Color("240")).Bold(true)

    tabSelected = lipgloss.NewStyle().
        Padding(0, 2).
        Foreground(lipgloss.Color("230")).
        Background(lipgloss.Color("#FF69B4")). // <- removed trailing dot
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
            renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
        } else {
            renderedTabs = append(renderedTabs, inactiveTabStyle.Render(t))
        }
    }
    return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func renderTodos(todos *todo.Todos) string {
    if len(*todos) == 0 {
        return lipgloss.NewStyle().
            Italic(true).
            Foreground(lipgloss.Color("238")).
            Render("No todos yet!")
    }

    var rendered string
    for _, t := range *todos {
        taskLine := lipgloss.NewStyle().Bold(true).Render(t.Task)
        details := lipgloss.JoinHorizontal(
            lipgloss.Top,
            lipgloss.NewStyle().Width(20).Render("Topic: "+t.Topic),
            lipgloss.NewStyle().Width(20).Render("Status: "+string(t.Status)),
            lipgloss.NewStyle().Width(20).Render("Priority: "+string(t.Priority)),
        )
        content := lipgloss.JoinVertical(lipgloss.Left, taskLine, details)
        rendered += todoBox.Render(content) + "\n\n"
    }
    return rendered
}




func renderProgress(step step, totalSteps int) string {
    progress := (int(step) + 1) * 100 / totalSteps
    filled := strings.Repeat("■", progress/10)
    empty := strings.Repeat("□", 10-(progress/10))

    barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
    emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

    return fmt.Sprintf("%s%s %d%%", barStyle.Render(filled), emptyStyle.Render(empty), progress)
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}


type model struct {
	step step
	totalSteps int
	taskInput textinput.Model
	progress progress.Model

	selectedTopic    int     //so these are like indexes to the string array
	selectedStatus   int
	selectedPriority int

	topics     []string
	statuses   []string
	priorities []string

	tabs      []string
    activeTab int

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
	var todos todo.Todos
    _ = todos.Load("todos.json")

	ti := textinput.New()
    ti.Placeholder = "What are we doing today?"
    ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true)
    ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
    ti.Focus()                              //great additional functionality
	p := progress.New(progress.WithDefaultGradient())
    p.Width = 50

	return model{
		step: stepTask,
		totalSteps: 5,
		taskInput: ti,
		progress: p,
		topics: []string{"Work", "Personal", "Health", "Study"}, //reminder- I have to change code to allow user topics
		statuses:   []string{"Pending", "In-progress", "Done"},
		priorities: []string{"Low", "Medium", "High"},

		tabs: []string{"Add", "List"},
		activeTab: 0,

		todos: &todos,
	}
}


func (m model) Init() tea.Cmd {   //Init can return a Cmd that could perform some initial I/O.
	return tea.Batch(
		textinput.Blink,
		tickCmd(), // start ticking right away
	) //blink for input

}


//The update function is called when ”things happen.” Its job is to look at what has happened and return an updated model in response. It can also return a Cmd to make more things happen

//The “something happened” comes in the form of a Msg, which can be any type. Messages are the result of some I/O that took place, such as a keypress, timer tick, or a response from a server.



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tickMsg:
        percent := float64(m.step+1) / float64(m.totalSteps)
        cmd := m.progress.SetPercent(percent)
        return m, tea.Batch(tickCmd(), cmd)
    case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit
		
		case "enter":
			if m.step < stepDone && m.step > stepList {
				m.step++
			} else if m.step == stepDone {
				// Save todo
				task := m.taskInput.Value()
				status := parseStatus(m.statuses[m.selectedStatus])
				priority := parsePriority(m.priorities[m.selectedPriority])
				topic := m.topics[m.selectedTopic]

				m.todos.Add(task, status, priority, topic)

				_ = m.todos.Store(".todos.json")

				// Go back to list view
				m.step = stepList
				m.taskInput.SetValue("") // reset
			}
		case "left":
			switch m.step {
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

		case "right":
			switch m.step {
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
		case ">":
			if m.step < stepDone {
				m.step++
			}
			return m, nil

		case "<":
			if m.step > stepList {
				m.step--
			}
			return m, nil
			}
	}
	// Task input only updates on first step
	if m.step == stepTask {
		var cmd tea.Cmd
		m.taskInput, cmd = m.taskInput.Update(msg)
		return m, cmd
	}
		
	return m,nil
}

//We look at the model in its current state and use it to return a string. That string is our UI!


func (m model) View() string {
	s := renderTabs(m) + "\n\n"

	switch m.step {

	case stepList:
		s += lipgloss.NewStyle().Bold(true).Render("Your Todos:\n")
		s += renderTodos(m.todos)
	case stepTask:
		s += "What are we doing today?\n" + m.taskInput.View()
		

	case stepTopic:
		s += "Choose a topic:\n"
		for i, t := range m.topics {
			if i == m.selectedTopic {
				s += selectedItem.Render(t)
			} else {
				s += itemStyle.Render(t)
			}
		}
	case stepStatus:
		s += "Pick a status:\n"
		for i, st := range m.statuses {
			if i == m.selectedStatus {
				s += selectedItem.Render(st)
			} else {
				s += itemStyle.Render(st)
			}
		}

	case stepPriority:
		s += "Select priority:\n"
		for i, p := range m.priorities {
			if i == m.selectedPriority {
				s += selectedItem.Render(p)
			} else {
				s += itemStyle.Render(p)
			}
		}
	
	case stepDone:
		s += "Okay, on it now!\n\n"
		s += fmt.Sprintf("Task: %s\n", m.taskInput.Value())
		s += fmt.Sprintf("Topic: %s\n", m.topics[m.selectedTopic])
		s += fmt.Sprintf("Status: %s\n", m.statuses[m.selectedStatus])
		s += fmt.Sprintf("Priority: %s\n", m.priorities[m.selectedPriority])
		s += "\nPress Enter to save & quit.\n"
	}
	s += "\n\nUse '<' and '>' to move tabs, 'n' to add a task, 'q' to quit.\n"
	s += "\n" + m.progress.View() + "\n"

	return s
}

// StartApp is the entrypoint from cmd/main.go
func StartApp() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}