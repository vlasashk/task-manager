package todo

type Repo interface {
	CreateTask(task Task) (Task, error)
	DeleteTaskByID(taskID string) error
	GetTaskByID(taskID string) (Task, error)
}
