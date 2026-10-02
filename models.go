package main

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// corpo esperado no POST
type titleReceiver struct {
	Title string `json:"title"`
}

// corpo esperado no PUT
type updatedTask struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
