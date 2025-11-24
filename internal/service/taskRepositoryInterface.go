package service

// интерфейс для репозитория задач (служит связующим звеном между двумя структурами: 
// реальной структурой TaskRepo и мок-структурой MockTaskRepository)

// то есть TaskRepoInterface описывает контракт, 
// который должен быть реализован любым объектом, претендующим на роль репозитория
type TaskRepoInterface interface {
	Create(text string, done *bool) (*TaskStruct, error)
	GetAll(tasks *[]TaskStruct) error
	GetByID(task *TaskStruct, id uint) error
	Update(task *TaskStruct) error
	Delete(task *TaskStruct, id uint) error 
}