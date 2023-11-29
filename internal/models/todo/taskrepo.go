package todo

type Repo interface {
	CreateTask(task TaskReq) (Task, error)
	DeleteTask(taskID string) error
	GetTask(taskID string) (Task, error)
	UpdateTask(task TaskReq, taskID string) (Task, error)
}
