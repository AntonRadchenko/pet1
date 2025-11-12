package service

// структура хранилища тасок
type TaskStruct struct {
	ID   uint `gorm:"primaryKey;autoIncrement"` 
	Task string `json:"task"`
}

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
}

// единый формат ошибок
type ErrorStruct struct {
	Error string `json:"error"`
}