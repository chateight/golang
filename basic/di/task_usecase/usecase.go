package task_usecase

import "di/task"

// Saveメソッドの実装はTaskRepositoryで行っている
type TaskRepositoryInterface interface {
	Save(task.Task) (task.Task, error)
}

// structのフィールドにTaskRepositoryInterfaceを持たせる
// こうすることで、CreateTaskメソッドでDBへの保存処理を呼ぶことができる
type TaskUsecase struct {
	repo TaskRepositoryInterface
}

func NewTaskUsecase(repo TaskRepositoryInterface) TaskUsecase {
	return TaskUsecase{repo: repo}
}

func (usecase TaskUsecase) CreateTask(title string) (task.Task, error) {
	t := task.Task{Title: title}
	task, err := usecase.repo.Save(t)
	return task, err
}