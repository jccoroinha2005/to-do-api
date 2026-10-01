package main

var tasks = []Task{
	{ID: 1, Title: "Update project documentation", Done: false},
	{ID: 2, Title: "Deploy staging environment", Done: true},
	{ID: 3, Title: "Review pull request", Done: true},
	{ID: 4, Title: "Schedule dentist appointment", Done: false},
}

var nextID = len(tasks) + 1 //proximo ID livre; so aumenta, entao nao repete ID depois de um DELETE

func findTaskIndex(id int) int {
	for i := range tasks {
		if tasks[i].ID == id {
			return i
		}
	}
	return -1
}
