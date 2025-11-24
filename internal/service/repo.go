package service

import (
	"strings"
	"time"

	"github.com/AntonRadchenko/WebPet1/internal/db"
)

// 2. repo-слой (руки)

// TaskRepo — слой доступа к данным
// здесь хранятся методы, которые напрямую обращаются к бд
// этот слой ничего не думает, не валидирует — просто делает CRUD-запросы

type TaskRepo struct{}

// Create - добавляет новую задачу в таблицу
func (r *TaskRepo) Create(text string, done *bool) (*TaskStruct, error) {
	// создаем сущность 
	task := &TaskStruct{
		Task: text,
		IsDone: *done,
	}

	if done != nil { // если done не пустой, то есть был передан в бади
		task.IsDone = *done // то обновляем его по указателю
	}

	err := db.DB.Create(task).Error
	if err != nil {
		return nil, err
	}
	return task, err // передаем объект задачи обратно в сервис
}

// GetAll - возвращает все задачи из таблицы
func (r *TaskRepo) GetAll(tasks *[]TaskStruct) error {
	err := db.DB.Find(&tasks).Error
	if err != nil && strings.Contains(err.Error(), "relation") {
		// если таблицы нет - возвращаем пустой массив []
		*tasks = []TaskStruct{} 
		return nil
	}
	return err
}

// GetByID - возвращает задачу по ID
func (r *TaskRepo) GetByID(task *TaskStruct, id uint) error {
	return db.DB.First(&task, "id = ?", id).Error
}

// Update - обновляет задачу (текст задачи)
func (r *TaskRepo) Update(task *TaskStruct) error {
	task.UpdatedAt = time.Now()
	return db.DB.Save(&task).Error
}

// Delete - удаляет задачу по ID
func (r *TaskRepo) Delete(task *TaskStruct, id uint) error {
	now := time.Now()
	task.DeletedAt = &now
	return db.DB.Delete(&task, id).Error
}