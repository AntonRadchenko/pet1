package taskService

import "time"

// модель базы данных 
type TaskStruct struct {
	ID   uint `gorm:"primaryKey;autoIncrement"` 
	Task string 
	IsDone bool 
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt *time.Time 
}

// единый формат ошибок
// type ErrorStruct struct {
// 	Error string `json:"error"`
// }