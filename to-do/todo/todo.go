package todo

import (
    "encoding/json"
    "errors"
    "os"
    "time"
)


type Status string
const (
	Pending Status = "Pending?!"
	InProgress Status = "In-progress..."
	Done Status = "DONE!"
)


type Priority int    //it’s still an int, That way, you can’t accidentally assign a random int where a Priority is expected, unless you explicitly cast.
const (
    Low Priority = iota   //iota is a Go keyword that auto-increments with each constant inside the const block.
    Medium
    High
)


type item struct {
	Task string
	Status Status
	Priority Priority
	Topic string
	Created time.Time
	Completed *time.Time   // pointer so it can be nil if not completed
	TimeWorked time.Duration

}

type Todos []item

//(t *Todos) is called the receiver. this function (Add) belongs to the type Todos. Inside the function, t acts like a variable that represents the Todos instance you’re working on. It’s very similar to self

func (t *Todos) Add(task string, status Status, priority Priority, topic string) {
	todo := item{
		Task: task,
		Status: status,
		Priority: priority,
		Topic: topic,
		Created: time.Now(),
		Completed: nil,
		TimeWorked: 0,
	}
	*t = append(*t, todo)
}


func (t *Todos) Complete(index int) error {
	ls := *t 
	now := time.Now()
	if index<=0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].Completed = &now
	ls[index-1].Status = Done
	ls[index-1].TimeWorked = time.Since(ls[index-1].Created)
	return nil 
}


func (t *Todos) Delete(index int) error {
	ls := *t
	if index<=0 || index > len(ls) {
		return errors.New("invalid index")
	}

	newTodos := make([]item, 0, len(ls))
    for _, todo := range ls {
        if !(todo.Status == Done && time.Now().Sub(*todo.Completed) > 24*time.Hour) {
            newTodos = append(newTodos, todo)
        }
    }


	newTodos = append(ls[:index-1],ls[index:]...) //slice to have everything - current index 
	*t = newTodos
	return nil
}


func (t *Todos) Edit(index int, task string, status string, priority int, topic string) error {
	ls := *t
	now := time.Now()
	if index<=0 || index > len(ls) {
		return errors.New("invalid index")
	}
	todo := &ls[index-1]
	todo.Task = task
    todo.Status = Status(status)
    todo.Priority = Priority(priority)
    todo.Topic = topic

	// If status is Done and Completed not yet set, set it to now
	if Status(status) == Done && todo.Completed == nil {
        todo.Completed = &now
        todo.TimeWorked = time.Since(todo.Created)
    }

	return nil
}


func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {                            //to catch error but why tf is there no try catch block in go? or inbuilt ones?
		if errors.Is(err, os.ErrNotExist) {
			// no file yet → treat as empty todo list
			return nil
		}
		return err
	}

	if len(file) == 0 {           
		return nil
	}
	err = json.Unmarshal(file, t)             //unmarshaling file Input: file → a []byte slice that contains JSON text (the serialized form of your todos). Output: t → your Go value (*Todos), which is a struct/slice in memory
	if err != nil {
		return err 
	}

	return err
}

func (t *Todos) Store(filename string) error {
		data, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			return err
		}
	
		return os.WriteFile(filename, data, 0644) //0644 is permissions 
}
