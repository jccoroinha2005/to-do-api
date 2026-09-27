package main

import (
	"encoding/json" //transformar dados de Golang para JSON
	"log"           //mensagens no terminal sobre o servidor
	"net/http"      //criar servidor e lidar com requisiçoes
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Update project documentation", Done: false},
	{ID: 2, Title: "Deploy staging environment", Done: true},
	{ID: 3, Title: "Review pull request", Done: true},
	{ID: 4, Title: "Schedule dentist appointment", Done: false},
}

type titleReceiver struct {
	Title string `json:"title"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json") //se o metodo for igual a GET, set um header
		json.NewEncoder(w).Encode(tasks)                   // crie um encoder e transforme a variavel tasks em JSON
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		var t titleReceiver
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "invalid request body!", http.StatusBadRequest)
			return
		}
		newID := len(tasks) + 1
		newTask := Task{ID: newID, Title: t.Title, Done: false}
		tasks = append(tasks, newTask)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)
	default:
		http.Error(w, "method not allowed!", http.StatusMethodNotAllowed) //imprima method not allowed e torne o valor do erro = 405
		return
	}
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	log.Println("servidor rodando em http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
