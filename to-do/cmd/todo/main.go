package main

import (
	"fmt"
	"flag"
	"os"
	// "github.com/spf13/cobra"
	"to-do/todos"

)
const (
	todoFile = ".todos.json"
)
func main() {
	add := flag.Bool("add", false,"Add a new todo")

	flag.Parse()
	todos := &todo.Todos{}

	if err := todo.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1);
	}

	switch{
	case *add: 
		todos.Add("Sample todo", Pending, 1, "Trial")
	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(1);
	}

}