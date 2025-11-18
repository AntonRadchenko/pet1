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
func (r *TaskRepo) Create(task *TaskStruct) error {
	return db.DB.Create(task).Error
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

