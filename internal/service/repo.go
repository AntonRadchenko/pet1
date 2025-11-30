package service

import (
	"strings"
	"time"

	"github.com/AntonRadchenko/WebPet1/internal/db"
	"gorm.io/gorm"
)

// 2. repo-слой (руки)

// TaskRepo — слой доступа к данным
// здесь хранятся методы, которые напрямую обращаются к бд
// этот слой ничего не думает, не валидирует — просто делает CRUD-запросы

// интерфейс для репозитория задач (служит связующим звеном между двумя структурами:
// реальной структурой TaskRepo и мок-структурой MockTaskRepository);
// То есть TaskRepoInterface описывает контракт,
// который должен быть реализован любым объектом, претендующим на роль репозитория
type TaskRepoInterface interface {
	Create(task *TaskStruct) (*TaskStruct, error) // исправлена пока только сигнатура этого метода
	GetAll() ([]TaskStruct, error)
	GetByID(id uint) (TaskStruct, error)
	Update(task *TaskStruct) (*TaskStruct, error)
	Delete(task *TaskStruct, id uint) error
}

type TaskRepo struct{}

// Create - добавляет новую задачу в таблицу
func (r *TaskRepo) Create(task *TaskStruct) (*TaskStruct, error) {
	err := db.DB.Create(task).Error // передаем указатель в ORM
	if err != nil {
		return nil, err
	}
	return task, nil // передаем объект задачи обратно в сервис
}

// GetAll - возвращает все задачи из таблицы
func (r *TaskRepo) GetAll() ([]TaskStruct, error) {
	var tasks []TaskStruct

	err := db.DB.Find(&tasks).Error
	if err != nil  {
		if strings.Contains(err.Error(), "relation") {
		// если таблицы нет - возвращаем пустой массив []
		return []TaskStruct{}, nil 
		}
		// возвращаем ошибку, если эта ошибка не изза отсутствия таблицы
		return nil, err
	}
	return tasks, nil
}

// GetByID - возвращает задачу по ID
func (r *TaskRepo) GetByID(id uint) (TaskStruct, error) {
	var task TaskStruct
	err := db.DB.First(&task, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// если задачи нет,
			return TaskStruct{}, nil
		}
		return task, err 
	}
	return task, nil
}

// Update - обновляет задачу (текст задачи)
func (r *TaskRepo) Update(task *TaskStruct) (*TaskStruct, error) {
	task.UpdatedAt = time.Now()
	err := db.DB.Save(&task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete - удаляет задачу по ID
func (r *TaskRepo) Delete(task *TaskStruct, id uint) error {
	now := time.Now()
	task.DeletedAt = &now
	err := db.DB.Delete(task).Error
	if err != nil {
		return err
	}
	return nil
}
