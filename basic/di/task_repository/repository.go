package task_repository

import "di/task"

type TaskRepository struct{}

// TaskRepositoryInterfaceを満たすようにSaveメソッドを実装
func (repo TaskRepository) Save(t task.Task) (task.Task, error) {
	// 実際にはDBへの保存処理を行う
	task := task.Task{
		ID: 1,
		Title: t.Title,
	}
	return task, nil
}
