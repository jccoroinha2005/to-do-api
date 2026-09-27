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

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { //se o metodo da requisiçao nao for igual ao metodo GET
		http.Error(w, "method not allowed!", http.StatusMethodNotAllowed) //imprima method not allowed e torne o valor do erro = 405
		return
	}

	w.Header().Set("Content-Type", "application/json") //se o metodo for igual a GET, set um header
	json.NewEncoder(w).Encode(tasks)                   // crie um encoder e transforme a variavel tasks em JSON
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	log.Println("servidor rodando em http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
