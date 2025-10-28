package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// глобальная переменная
var task string = "Anton"

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода 405
	}
	var requestBody RequestBody
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		fmt.Println("err: ", err)
	}
	// записываем содержимое в нашу переменную task
	task = requestBody.Task
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода 405
	}
	fmt.Println("Hello,", task) // пишем в терминал
	_, err := w.Write([]byte("hello, " + task + "\n")) // пишем ответ пользователю
	if err != nil {
		fmt.Println("fail to write HTTP response: ", err)
	}
}

func main() {
	http.HandleFunc("/hello", GetHandler)
	http.HandleFunc("/task", PostHandler)
	if err := http.ListenAndServe(":9092", nil); err != nil { // слушаем порт 9092
		fmt.Println("Ошибка во время работы HTTP сервера: ", err)
	}
}