package main

import (
	"context"
	"log"      //mensagens no terminal sobre o servidor
	"net/http" //criar servidor e lidar com requisiçoes
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	db = pool

	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/{id}", idTaskHandler)

	log.Println("server running on http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
