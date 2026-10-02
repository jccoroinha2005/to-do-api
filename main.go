package main

import (
	"log"      //mensagens no terminal sobre o servidor
	"net/http" //criar servidor e lidar com requisiçoes
)

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/{id}", idTaskHandler)

	log.Println("server running on http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
