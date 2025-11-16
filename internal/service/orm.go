package service

import "time"

// модель базы данных 
type TaskStruct struct {
	ID   uint `gorm:"primaryKey;autoIncrement"` 
	Task string `json:"task"`
	IsDone bool `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// структура тела запроса
// type RequestBody struct {
// 	Task string `json:"task"`
// 	IsDone *bool `json:"is_done"`
// }

// единый формат ошибок
type ErrorStruct struct {
	Error string `json:"error"`
}