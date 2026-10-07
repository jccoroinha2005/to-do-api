# to-do-api

RESTful To-Do API built in Go, with full CRUD operations, backed by PostgreSQL and packaged with Docker.

## Tech stack

- **Go**: standard library `net/http`, no web framework
- **PostgreSQL 18**: data storage
- **pgx v5** (`pgxpool`): PostgreSQL driver with connection pooling
- **Docker & Docker Compose**: run the API and the database together with one command

## Getting started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) with Docker Compose

You do **not** need Go or PostgreSQL installed. Docker creates a private database on your own machine, so nothing is shared with anyone else.

### Run it

1. Clone the repository:

   ```bash
   git clone git@github.com:jccoroinha2005/to-do-api.git
   cd to-do-api
   ```

2. Create your `.env` file from the example and change the password:

   ```bash
   cp .env.example .env
   ```

   Open `.env` and replace `change_me` with any password you like. It is only used by the database created on your machine. Use the same password in both lines.

3. Start everything:

   ```bash
   docker compose up --build -d
   ```

   On the first run, the database is created and `schema.sql` creates the `tasks` table with 4 example tasks.

4. Check that it works:

   ```bash
   curl localhost:8080/tasks
   ```

The API is now running at `http://localhost:8080`.

### Stop it

```bash
docker compose down        # stops the containers, keeps your data
docker compose down -v     # stops the containers and deletes all data
```

Use `down -v` to start again from scratch (the next run recreates the table with the 4 example tasks).

## API endpoints

| Method | Endpoint      | Description        | Body                                  |
|--------|---------------|--------------------|---------------------------------------|
| GET    | `/tasks`      | List all tasks     |                                       |
| POST   | `/tasks`      | Create a task      | `{"title": "..."}`                    |
| GET    | `/tasks/{id}` | Get one task       |                                       |
| PUT    | `/tasks/{id}` | Update a task      | `{"title": "...", "done": true}`      |
| DELETE | `/tasks/{id}` | Delete a task      |                                       |

A task looks like this:

```json
{
  "id": 1,
  "title": "Update project documentation",
  "done": false
}
```

### Examples

```bash
# List all tasks
curl localhost:8080/tasks

# Create a task
curl -X POST -d '{"title":"Buy groceries"}' localhost:8080/tasks

# Get task 5
curl localhost:8080/tasks/5

# Update task 5
curl -X PUT -d '{"title":"Buy groceries","done":true}' localhost:8080/tasks/5

# Delete task 5
curl -X DELETE localhost:8080/tasks/5
```

Add `-i` to any `curl` command to see the HTTP status code.

### Status codes

| Code | Meaning                                                      |
|------|--------------------------------------------------------------|
| 200  | Success (GET, PUT)                                           |
| 201  | Task created (POST)                                          |
| 204  | Task deleted (DELETE), no body                               |
| 400  | Invalid ID, invalid JSON, or empty title                     |
| 404  | Task not found                                               |
| 405  | Method not allowed                                           |
| 500  | Database error                                               |

## Project structure

```
.
├── main.go              # loads config, connects to the database, starts the server
├── handlers.go          # HTTP handlers for /tasks and /tasks/{id}
├── models.go            # Task struct and request bodies
├── store.go             # shared database connection pool
├── schema.sql           # creates the tasks table and example data
├── Dockerfile           # multi-stage build of the API image
├── docker-compose.yml   # API + PostgreSQL
└── .env.example         # template for your .env file
```

## Running without Docker

If you prefer, you can run the API directly. You need Go 1.26+ and a running PostgreSQL.

1. Create a database and a user, then run `schema.sql` on it.
2. Create your `.env` from `.env.example` and set `DATABASE_URL` to your own database.
3. Run:

   ```bash
   go run .
   ```

## My choices and experience

- **Starting from scratch in a single file:** I began with everything in one `main.go`: the imports, the structs, the in-memory storage, all the handler functions and the `main` function that starts the server.
- **One handler at first:** the API started with a single handler, `tasksHandler`, serving one endpoint (`/tasks`) with two methods, GET and POST. Those two don't need an ID, so one function was enough.
- **A second handler for the ID routes:** for PUT and DELETE I needed to act on one specific task, so I wrote `idTaskHandler`, which reads the ID from the URL (`/tasks/{id}`). I also added a GET by ID to fetch a single task.
- **Splitting the code into files:** to keep the logic organized and avoid mixing responsibilities, I split `main.go` into separate files, a common practice in larger codebases:
  - `handlers.go`: the endpoint logic and all the CRUD methods, the most important part of the API.
  - `store.go`: defines the global `db` variable (short for database), used by the handlers to talk to the database.
  - `models.go`: defines the structure of a task and the request bodies.
  - `main.go`: the setup work: reading the `.env`, connecting to the database, checking that it is reachable, registering the handlers and starting the server.
- **Moving from memory to PostgreSQL:** the first version stored the tasks in a slice, so every task was lost when the server stopped. I migrated to PostgreSQL, a popular and reliable relational database, which solves exactly that problem.
- **pgx and parameterized queries:** I used [pgx](https://github.com/jackc/pgx), the PostgreSQL driver for Go, with its `pgxpool` package, which manages a pool of connections. The queries use parameters (`$1`, `$2`) instead of concatenating strings.
- **Docker to make it easy to run:** I containerized the project with Docker and Docker Compose, so anyone who wants to run and play with the code only needs a few commands in the terminal.
- **Keeping the old version:** instead of deleting the in-memory version, I tagged it as `v1.0`, so both versions can be read side by side in the repository history.

## Versions

- **`v2.0`**: PostgreSQL + Docker (current version, on `main`)
- **`v1.0`**: in-memory version, without a database. To read that code:

  ```bash
  git checkout v1.0
  git switch main   # to come back
  ```

  Or choose the `v1.0` tag in the branch selector on GitHub.

## Concepts (for the curious)

A short glossary of the ideas behind this project, in case you want to dig deeper.

### The API

- **API** (Application Programming Interface): a set of rules that lets one piece of software talk to another. Think of a waiter in a restaurant: you (the client) tell the waiter what you want, the waiter takes it to the kitchen (the server) and brings back the result. You never need to know how the kitchen works.
- **REST** (Representational State Transfer): a style for designing web APIs on top of HTTP. This project is a RESTful API, which means:
  - Things are called **resources** and have their own URL (`/tasks`, `/tasks/{id}`).
  - The **HTTP method** says what to do with the resource (GET, POST, PUT, DELETE).
  - It is **stateless**: each request carries everything the server needs, and the server does not remember previous requests.
  - Data travels as **JSON** and the result is reported with an HTTP **status code**.
- **CRUD**: the four basic operations on data: **C**reate, **R**ead, **U**pdate and **D**elete. Each one maps to an HTTP method here, and to an SQL command in the database:

  | CRUD   | HTTP method | SQL      |
  |--------|-------------|----------|
  | Create | POST        | `INSERT` |
  | Read   | GET         | `SELECT` |
  | Update | PUT         | `UPDATE` |
  | Delete | DELETE      | `DELETE` |

- **Endpoint**: a URL together with an HTTP method, for example `GET /tasks`.
- **JSON** (JavaScript Object Notation): a simple text format for data, like `{"id": 1, "title": "...", "done": false}`. Almost every web API uses it.
- **Status code**: a number in the response that tells how the request went. `2xx` means success, `4xx` means the client made a mistake (such as `404` not found) and `5xx` means the server failed.

### The database

- **Relational database**: stores data in tables made of rows and columns. PostgreSQL is one of the most popular open-source ones.
- **SQL** (Structured Query Language): the language used to talk to a relational database. `schema.sql` in this repository is written in SQL.
- **Primary key and `SERIAL`**: the primary key (`id`) identifies each row. `SERIAL` makes the database generate it automatically, always increasing. That is why a deleted id is never reused.
- **Driver**: the library that lets a program talk to a specific database. Here it is `pgx`.
- **Connection pool**: a set of database connections kept open and reused by the requests, instead of opening a new one every time. That is what `pgxpool` does.
- **SQL injection**: an attack where someone sends SQL code disguised as data. Parameterized queries (`$1`, `$2`) prevent it, because the values are sent separately from the SQL text.

### Configuration and Docker

- **Environment variables and `.env`**: configuration kept outside the code, so secrets like passwords never get committed. The `.env` file is ignored by git, and `.env.example` shows what it should look like.
- **Docker**: packages an application and everything it needs to run into a **container**, so it behaves the same on any machine.
- **Image vs. container**: an image is the recipe (built from the `Dockerfile`) and a container is a running instance of it.
- **Multi-stage build**: the `Dockerfile` first compiles the Go code in a big image with the Go toolchain, then copies only the finished binary into a small final image.
- **Docker Compose**: describes several containers in one file (`docker-compose.yml`) and starts them together, here the API and PostgreSQL.
- **Volume**: storage that lives outside the container, so the database data survives when the containers are removed. `docker compose down -v` deletes it.

### Git and GitHub

- **Git**: a version control system created by Linus Torvalds, the creator of the Linux kernel. I used it to version and save all of my code, so every commit I made to each file can be seen on GitHub.
- **Tag**: a name attached to a specific commit, used to mark versions. Here, `v1.0` and `v2.0`.
