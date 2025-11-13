package service

import "time"

// структура хранилища тасок
type TaskStruct struct {
	ID   uint `gorm:"primaryKey;autoIncrement"` 
	Task string `json:"task"`
	IsDone bool `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	// в DeletedAt указатель нужен,
	// чтобы поле могло быть пустым (nil),
	// а не содержать фальшивую дату 0001-01-01.
}

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
	IsDone *bool `json:"is_done"`
}

// единый формат ошибок
type ErrorStruct struct {
	Error string `json:"error"`
}