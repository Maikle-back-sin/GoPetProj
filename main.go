package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp/v3"
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

	err = json.NewDecoder(file).Decode(&tasks)
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

func getTaskByID(id int) ([]Task, *Task, error) {
	tasks, err := loadTask()
	if err != nil {
		return nil, nil, err
	}
	for i := range tasks {
		if tasks[i].ID == id {
			return tasks, &tasks[i], nil
		}
	}
	return tasks, nil, fmt.Errorf("task with ID %d not found", id)
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
		Status:      "To-Do",
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

func updateTask(id int, newDesc string) error {
	tasks, task, err := getTaskByID(id)
	if err != nil {
		return err
	}
	task.Description = newDesc
	task.UpdatedAt = time.Now()
	return saveTask(tasks)
}

func updateTaskStatusInProgress(id int) error {
	tasks, task, err := getTaskByID(id)
	if err != nil {
		return err
	}
	task.Status = "In-Progress"
	task.UpdatedAt = time.Now()
	return saveTask(tasks)
}

func updateTaskStatusDone(id int) error {
	tasks, task, err := getTaskByID(id)
	if err != nil {
		return err
	}
	task.Status = "Done"
	task.UpdatedAt = time.Now()
	return saveTask(tasks)
}

func updateTaskStatusInToDO(id int) error {
	tasks, task, err := getTaskByID(id)
	if err != nil {
		return err
	}
	task.Status = "To-Do"
	task.UpdatedAt = time.Now()
	return saveTask(tasks)
}

func deleteTask(id int) error {
	tasks, err := loadTask()
	if err != nil {
		return err
	}

	found := false
	newTasks := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue // пропускаем задачу, которую нужно удалить
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	err = saveTask(newTasks)
	if err != nil {
		return err
	}

	fmt.Println("Deleted task with ID:", id)
	return nil
}

func getAllTasksWithStatus(status string) ([]Task, error) {
	resTasks := []Task{}
	tasks, err := loadTask()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.Status == status {
			resTasks = append(resTasks, task)
		}
	}

	return resTasks, err
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

		// Парсим команду и аргументы для add
		if strings.HasPrefix(input, "add") {
			desc := strings.TrimSpace(strings.TrimPrefix(input, "add"))

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

		// Парсим команду и аргументы для update
		if strings.HasPrefix(input, "update") {
			parts := strings.SplitN(input, " ", 3)

			if len(parts) < 3 {
				fmt.Println("Usage: update <id> <new description>")
				continue
			}

			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid task ID.", parts[1])
				continue
			}

			newDesc := parts[2]
			err = updateTask(id, newDesc)
			if err != nil {
				fmt.Println("Error updating task:", err)
			} else {
				fmt.Println("Updated task:", newDesc)
			}
			continue
		}

		// Парсим команду и аргументы для delete
		if strings.HasPrefix(input, "delete") {
			fields := strings.Fields(input)
			if len(fields) < 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			id, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Invalid task ID.", fields[1])
				continue
			}
			err = deleteTask(id)
			if err != nil {
				fmt.Println("Error deleting task:", err)
			}
			continue
		}

		// Парсим команду и аргументы для mark-to-do
		if strings.HasPrefix(input, "mark-to-do") {
			parts := strings.Fields(input)

			if len(parts) < 2 {
				fmt.Println("Usage: update <id> <new description>")
				continue
			}

			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid task ID.", parts[1])
				continue
			}

			err = updateTaskStatusInToDO(id)
			if err != nil {
				fmt.Println("Error updating task-status:", err)
			} else {
				fmt.Println("Updated task-status to mark-to-do seccessfully:")
			}
			continue
		}

		// Парсим команду и аргументы для mark-in-progress
		if strings.HasPrefix(input, "mark-in-progress") {
			parts := strings.Fields(input)

			if len(parts) < 2 {
				fmt.Println("Usage: update <id> <new description>")
				continue
			}

			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid task ID.", parts[1])
				continue
			}

			err = updateTaskStatusInProgress(id)
			if err != nil {
				fmt.Println("Error updating task-status:", err)
			} else {
				fmt.Println("Updated task-status to mark-in-progress seccessfully:")
			}
			continue
		}

		// Парсим команду и аргументы для mark-done
		if strings.HasPrefix(input, "mark-done") {
			parts := strings.Fields(input)

			if len(parts) < 2 {
				fmt.Println("Usage: update <id> <new description>")
				continue
			}

			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid task ID.", parts[1])
				continue
			}

			err = updateTaskStatusDone(id)
			if err != nil {
				fmt.Println("Error updating task-status:", err)
			} else {
				fmt.Println("Updated task-status to DONE seccessfully:")
			}
			continue
		}

		// Парсим команду и аргументы для list
		if strings.HasPrefix(input, "list") {
			parts := strings.Fields(input)
			if len(parts) < 2 {
				tasks, err := loadTask()
				if err != nil {
					fmt.Println("Error showing tasks:", err)
				} else {
					pp.Println(tasks)
				}
			} else {
				tasks, err := getAllTasksWithStatus(parts[1])
				if err != nil {
					fmt.Println("Error showing tasks:", err)
				} else {
					pp.Println(tasks)
				}
			}
			continue
		}

		fmt.Println("Unknown command:", input)
	}
}
