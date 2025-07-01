package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")

		if ok := scanner.Scan(); !ok {
			fmt.Println("Scan error.")
		}

		text := scanner.Text()

		fields := strings.Fields(text)

		if len(fields) == 0 {
			fmt.Println("Empty input.")
		}

		if fields[0] == "add" {
			fmt.Println("Your enter is add : ")
		}

		if fields[0] == "update" {
			fmt.Println("Your enter is update : ")
		}

		if fields[0] == "delete" {
			fmt.Println("Your enter is delete: ")
		}

		if fields[0] == "mark-todo" {
			fmt.Println("Your enter is mark-todo: ")
		}

		if fields[0] == "mark-in-progress" {
			fmt.Println("Your enter is mark-in-progress: ")
		}

		if fields[0] == "mark-done" {
			fmt.Println("Your enter is mark-done: ")
		}

		if fields[0] == "list" && len(fields) == 1 {
			fmt.Println("Default list of all tasks: ")
		}

		if fields[0] == "list" && len(fields) > 1 {
			if fields[1] == "done" {
				fmt.Println("Your enter list done: ")
			}

			if fields[1] == "todo" {
				fmt.Println("Your enter list todo: ")
			}

			if fields[1] == "in-progress" {
				fmt.Println("Your enter list in-progress: ")
			}

		}
	}
}

type Task struct {
	id          int
	description string
	status      string
	createdAt   time.Time
	updatedAt   time.Time
}
