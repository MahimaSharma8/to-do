package tea 

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/MahimaSharma8/to-do/todo"
)


type step int 
const (
	stepTask step = iota
	stepTopic
	stepStatus
	stepPriority
	stepDone
)


type model struct {
	step step
	taskInput textinput.Model

	selectedTopic    int     //so these are like indexes to the string array
	selectedStatus   int
	selectedPriority int

	topics     []string
	statuses   []string
	priorities []string


	todos *todo.Todos
}



func initialModel() model {
	ti := textinput.New()                       
	ti.Placeholder = "What are we doing today?"
	ti.Focus()                                     //great additional functionality

	return model{
		step: stepTask,
		taskInput: ti,

		topics: []string{"Work", "Personal", "Health", "Study"}, //reminder- I have to change code to allow user topics
		statuses:   []string{"Pending", "In-progress", "Done"},
		priorities: []string{"Low", "Medium", "High"},

		todos: &todo.Todos{},
	}
}


func (m model) Init tea.Cmd {   //Init can return a Cmd that could perform some initial I/O.
	return textinput.Blink //blink for input
}


//The update function is called when ”things happen.” Its job is to look at what has happened and return an updated model in response. It can also return a Cmd to make more things happen

//The “something happened” comes in the form of a Msg, which can be any type. Messages are the result of some I/O that took place, such as a keypress, timer tick, or a response from a server.



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
    case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit
		
		case "enter":
			if m.step < stepDone {
				m.step++
			} else {
				// Save todo
				task := m.taskInput.Value()
				status := m.statuses[m.selectedStatus]
				priority := m.priorities[m.selectedPriority]
				topic := m.topics[m.selectedTopic]

				m.todos.Add(task, todo.Status(status), todo.Priority(priority), topic)

				todos.Store(todoFile)
				return m, tea.Quit
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
	s := "Good Morning~"   //reminder- change code to change greetings based on time of the day

	// Tabs
	tabs := []string{"Task", "Topic", "Status", "Priority", "Done"}
	for i, t := range tabs {
		if i == int(m.step) {
			s += fmt.Sprintf("[ %s ] ", t) //??
		} else {
			s += fmt.Sprintf("  %s   ", t)
		}
	}
	s += "\n\n"   //?? 

	switch m.step {
	case stepTask:
		s += "What are we doing today?\n"
		s += m.taskInput.View() + "\n"

	case stepTopic:
		s += "Choose a topic:\n"
		for i, t := range m.topics {
			if i == m.selectedTopic {
				s += fmt.Sprintf("[ %s ] ", t)
			} else {
				s += fmt.Sprintf("  %s   ", t)
			}
		}
	case stepStatus:
		s += "Pick a status:\n"
		for i, st := range m.statuses {
			if i == m.selectedStatus {
				s += fmt.Sprintf("[ %s ] ", st)
			} else {
				s += fmt.Sprintf("  %s   ", st)
			}
		}

	case stepPriority:
		s += "Select priority:\n"
		for i, p := range m.priorities {
			if i == m.selectedPriority {
				s += fmt.Sprintf("[ %s ] ", p)
			} else {
				s += fmt.Sprintf("  %s   ", p)
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

	// Progress bar
	progress := (int(m.step) + 1) * 100 / 5
	s += fmt.Sprintf("\n\n[%s%s] %d%%",
		strings.Repeat("#", progress/10),
		strings.Repeat("-", 10-(progress/10)),
		progress,
	)

	return s
}

// StartApp is the entrypoint from cmd/main.go
func StartApp() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}