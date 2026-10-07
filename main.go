package main

import (
	"context"
	"log"      //mensagens no terminal sobre o servidor
	"net/http" //criar servidor e lidar com requisiçoes
	"os"       //ler variáveis de ambiente

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// lê o .env se ele existir; no Docker não há .env, as variáveis vêm do docker-compose
	godotenv.Load()

	// cria o pool de conexões (ainda não conecta de verdade)
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// testa a conexão: se o banco estiver fora do ar, a API nem inicia
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	db = pool //guarda o pool na variável global de store.go, usada pelos handlers

	http.HandleFunc("/tasks", tasksHandler)
	// o {id} é lido dentro do handler com r.PathValue("id")
	http.HandleFunc("/tasks/{id}", idTaskHandler)

	log.Println("server running on http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
