package main

import (
	"di/task_repository"
	"di/task_usecase"
	"fmt"
)

func main() {
	repo := task_repository.TaskRepository{}
	usecase := task_usecase.NewTaskUsecase(repo)
	task, _ := usecase.CreateTask("DIの勉強")
	fmt.Printf("ID: %d, Title: %s\n", task.ID, task.Title) // => ID: 1, Title: DIの勉強
}
