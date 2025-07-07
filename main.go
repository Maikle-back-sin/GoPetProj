package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const fileName = "task.json"

func loadTask() ([]Task, error) {
	var tasks []Task

	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder((file)).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTask(tasks []Task) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}

func generateID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func getTaskByID(id int) Task {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return Task{}
	}
	defer file.Close()

	var tasks []Task
	err = json.NewDecoder((file)).Decode(&tasks)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return Task{}
	}

	for _, task := range tasks {
		if task.ID == id {
			return task
		}
	}

	return Task{}
}

func updateTask(id int) Task {
	return getTaskByID(id)
}

func addTask(description string) error {
	tasks, err := loadTask()

	if err != nil {
		return err
	}

	now := time.Now()

	newTask := Task{
		ID:          generateID(tasks),
		Description: description,
		Status:      "to-do",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	err = saveTask(tasks)
	if err != nil {
		return err
	}

	fmt.Println("Added new task:", newTask.Description)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")

		if !scanner.Scan() {
			fmt.Println("Scan error.")
			continue
		}

		input := scanner.Text()
		if strings.TrimSpace(input) == "" {
			fmt.Println("Empty input.")
			continue
		}

		// Парсим команду и аргументы
		if strings.HasPrefix(input, "add") {
			desc := strings.TrimPrefix(input, "add")

			if desc == "" {
				fmt.Println("Description cannot be empty.")
				continue
			}
			err := addTask(desc)
			if err != nil {
				fmt.Println("Error adding task:", err)
			}
			continue
		}

		if strings.HasPrefix(input, "update") {
			fields := strings.Fields(input)

			if len(fields) < 2 {
				fmt.Println("Please provide a task description.")
				continue
			}

			idStr := fields[1]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid task ID.", idStr)
				continue
			}

			task := updateTask(id)
			fmt.Println(task.Description)
			continue
		}

		fmt.Println("Unknow command:", input)
	}
}
