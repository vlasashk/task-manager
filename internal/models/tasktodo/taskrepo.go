package tasktodo

type Repo interface {
	CreateTask(task Request) (Task, error)
	DeleteTask(taskID string) error
	GetTask(taskID string) (Task, error)
	ListTasks(page uint, date string, status string) ([]Task, error)
	UpdateTask(task Request, taskID string) (Task, error)
}
