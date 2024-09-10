package handler

import (
	"code/internal/model"
	"code/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func tasksView(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	repo := repository.NewTaskRepository()
	defer repo.CloseConnection()

	idStr := strings.TrimPrefix(path, "/tasks/")

	if idStr != "" {

		taskId, err := strconv.Atoi(idStr)

		if err != nil {
			log.Fatal(err)
		}

		taskById(writer, request, taskId, repo)

	} else {
		if request.Method == http.MethodGet {
			tasksSlice := repo.GetAllTasks()
			jsonData, err := json.Marshal(tasksSlice)

			if err != nil {
				log.Fatal(err)
			}

			writer.WriteHeader(http.StatusOK)
			writer.Write(jsonData)
		}

		if request.Method == http.MethodPost {
			var task *model.Task
			err := json.NewDecoder(request.Body).Decode(&task)
			if err != nil {
				log.Fatal(err)
			}
			err = repo.AddTask(task)
			if err == nil {
				writer.WriteHeader(http.StatusCreated)
			} else {
				writer.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}

func taskById(w http.ResponseWriter, r *http.Request, taskId int, repo *repository.TaskRepository) {
	if r.Method == http.MethodGet {
		task, err := repo.GetTaskById(taskId)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		jsonData, err := json.Marshal(task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
	if r.Method == http.MethodPut {
		var task *model.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		repo.UpdateTask(task, taskId)
	}
	if r.Method == http.MethodDelete {
		err := repo.DeleteTask(taskId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func StartServer() {
	http.HandleFunc("/tasks/", tasksView)
	http.ListenAndServe("localhost:8000", nil)
}
