package main

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type titleReceiver struct {
	Title string `json:"title"`
}

type updatedTask struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
