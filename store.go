package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// pool de conexões com o PostgreSQL: é criado na main e usado por todos os handlers
var db *pgxpool.Pool
