package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FileName         = "task.json"
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func loadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.Open(FileName)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	file, err := os.Create(FileName)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(tasks)
}

func nextID(tasks []Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func findTask(tasks []Task, id int) (*Task, int) {
	for i, t := range tasks {
		if t.ID == id {
			return &tasks[i], i
		}
	}
	return nil, -1
}

func cmdAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: add \"description\"")
	}
	description := strings.Join(args[1:], " ")
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	now := time.Now()
	task := Task{
		ID:          nextID(tasks),
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks = append(tasks, task)
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
	return nil
}

func cmdUpdate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: update <id> \"new description\"")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	task, _ := findTask(tasks, id)
	if task == nil {
		return fmt.Errorf("task %d not found", id)
	}
	task.Description = strings.Join(args[2:], " ")
	task.UpdatedAt = time.Now()
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Println("Task updated")
	return nil
}

func cmdDelete(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: delete <id>")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	_, idx := findTask(tasks, id)
	if idx == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	tasks = append(tasks[:idx], tasks[idx+1:]...)
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Println("Task deleted")
	return nil
}

func cmdSetStatus(args []string, status string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: mark-%s <id>", status)
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	task, _ := findTask(tasks, id)
	if task == nil {
		return fmt.Errorf("task %d not found", id)
	}
	task.Status = status
	task.UpdatedAt = time.Now()
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("Task %d marked as %s\n", id, status)
	return nil
}

func cmdList(args []string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	filter := ""
	if len(args) > 1 {
		filter = args[1]
	}
	for _, t := range tasks {
		if filter == "" || t.Status == filter {
			fmt.Printf("[%d] %s [%s] (created: %s, updated: %s)\n", t.ID, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04"), t.UpdatedAt.Format("2006-01-02 15:04"))
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli <command> [arguments]")
		return
	}
	cmd := os.Args[1]
	args := os.Args[1:]
	var err error
	switch cmd {
	case "add":
		err = cmdAdd(args)
	case "update":
		err = cmdUpdate(args)
	case "delete":
		err = cmdDelete(args)
	case "mark-todo":
		err = cmdSetStatus(args, StatusTodo)
	case "mark-in-progress":
		err = cmdSetStatus(args, StatusInProgress)
	case "mark-done":
		err = cmdSetStatus(args, StatusDone)
	case "list":
		err = cmdList(args)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
}
