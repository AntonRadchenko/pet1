package service

import "github.com/AntonRadchenko/WebPet1/internal/db"

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
	return db.DB.Find(&tasks).Error
}

// GetByID - возвращает задачу по ID
func (r *TaskRepo) GetByID(task *TaskStruct, id uint) error {
	return db.DB.First(&task, "id = ?", id).Error
}

// Update - обновляет задачу (текст задачи)
func (r *TaskRepo) Update(task *TaskStruct) error {
	return db.DB.Save(&task).Error
}

// Delete - удаляет задачу по ID
func (r *TaskRepo) Delete(task *TaskStruct, id uint) error {
	return db.DB.Delete(&task, id).Error
}

