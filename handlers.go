package main

import (
	"encoding/json" //transformar dados de Golang para JSON
	"net/http"      //criar servidor e lidar com requisiçoes
	"strconv"
)

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
		newTask := Task{ID: nextID, Title: t.Title, Done: false}
		nextID++
		tasks = append(tasks, newTask)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)
	default:
		http.Error(w, "method not allowed!", http.StatusMethodNotAllowed) //imprima method not allowed e torne o valor do erro = 405
		return
	}
}

func idTaskHandler(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, "write a valid ID please!", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		i := findTaskIndex(id)
		if i == -1 {
			http.Error(w, "sorry, there's no task with this ID!", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks[i])
	case http.MethodPut:
		i := findTaskIndex(id)
		if i == -1 {
			http.Error(w, "sorry, there's no task with this ID!", http.StatusNotFound)
			return
		}
		var u updatedTask
		err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid request body!", http.StatusBadRequest)
			return
		}
		if u.Title == "" {
			http.Error(w, "write a valid Title please!", http.StatusBadRequest)
			return
		}
		tasks[i].Title = u.Title
		tasks[i].Done = u.Done
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks[i])
	case http.MethodDelete:
		i := findTaskIndex(id)
		if i == -1 {
			http.Error(w, "sorry, there's no task with this ID!", http.StatusNotFound)
			return
		}
		tasks = append(tasks[:i], tasks[i+1:]...)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed!", http.StatusMethodNotAllowed)
		return
	}
}
