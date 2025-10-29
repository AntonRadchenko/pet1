package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// глобальная переменная
var task = []string{}

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода 405
		return
	}
	var requestBody RequestBody
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// записываем содержимое в нашу переменную task
	task = append(task, requestBody.Task)
	fmt.Println("Task created succesfully!")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task created successfully\n"))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода 405
		return
	}
	if len(task) == 0 {
		fmt.Println("Hello, nobody yet!")
		_, err := w.Write([]byte("hello nobody yet!\n"))
		if err != nil {
			fmt.Println("fail to write HTTP response: ", err)
		}
		return 
	}
	fmt.Println("Hello,", task[len(task)-1])                        // пишем в терминал
	_, err := w.Write([]byte("hello, " + task[len(task)-1] + "\n")) // пишем ответ пользователю
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