package main

import (
	"code/internal/repository"
	"fmt"
)

func main() {
	repo := repository.NewTaskRepository()
	defer repo.CloseConnection()
	err := repo.AddTask("asd", "asdasd", "asdasd")
	fmt.Println(err)
	//repo.AddTask("new_task", "such a wonderful task", "2023-10-05T14:48:00Z")
	//repo.UpdateTask(7, "updated title", "new discra", "2023-10-05T14:48:00Z")
	// repo.DeleteTask(10)
}
