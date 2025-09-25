package main

import (
	"fmt"
	"os"
	"github.com/MahimaSharma8/to-do/tea"
)

func main() {
	if err := tea.StartApp(); err != nil {
		fmt.Println("Error running app:", err)
		os.Exit(1)
	}
}
