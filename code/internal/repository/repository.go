package repository

import (
	"code/internal/model"
	"code/internal/pkg/logger"

	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var AppLogger = logger.CreateNewLogger()

type TaskRepositoryInterface interface {
	CloseConnection()
}

type TaskRepository struct {
	db *sql.DB
}

func check(err error) {
	if err != nil {
		AppLogger.Fatal(err)
	}
}

func NewTaskRepository() *TaskRepository {
	db, err := connectDB()
	check(err)
	return &TaskRepository{db: db}
}

func (t *TaskRepository) CloseConnection() {
	err := t.db.Close()
	check(err)
	AppLogger.Info("Database connection was closed successfully")
}

func (t *TaskRepository) GetAllTasks() []*model.Task {
	rows, err := t.db.Query("SELECT title, description FROM tasks;")

	check(err)

	tasks := make([]*model.Task, 0)

	for rows.Next() {
		task := &model.Task{}
		err = rows.Scan(&task.Title, &task.Description)
		check(err)
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		fmt.Println(*task)
	}

	return tasks
}

func (t *TaskRepository) GetTaskById(id int) (*model.Task, error) {
	row := t.db.QueryRow("SELECT id, title, description FROM tasks WHERE id = $1", id)
	task := &model.Task{}
	row.Scan(&task.Id, &task.Title, &task.Description)
	if task.Id != id {
		return task, errors.New("task not found")
	}
	return task, nil
}

func (t *TaskRepository) AddTask(title string, description string, dueDate string) error {
	nowTime := time.Now().UTC().Format(time.RFC3339)
	_, err := t.db.Exec("INSERT INTO tasks (title, description, duedate, created_at, updated_at) VALUES ($1,$2,$3,$4,$4)", title, description, dueDate, nowTime)
	check(err)
	return err
}

func (t *TaskRepository) UpdateTask(id int, title string, description string, dueDate string) error {
	nowTime := time.Now().UTC().Format(time.RFC3339)
	res, err := t.db.Exec("UPDATE tasks SET title = $1, description = $2, duedate = $3, updated_at = $4 WHERE id = $5", title, description, dueDate, nowTime, id)
	check(err)
	rowsCount, err := res.RowsAffected()
	check(err)

	if !(rowsCount > 0) {
		return err
	}

	return nil
}

func (t *TaskRepository) DeleteTask(id int) error {
	res, err := t.db.Exec("DELETE FROM tasks WHERE id = $1;", id)
	check(err)
	rowsCount, err := res.RowsAffected()
	check(err)

	if !(rowsCount > 0) {
		return err
	}

	return nil
}

func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=localhost port=5432 sslmode=disable", username, dbname, password)

	db, err := sql.Open("postgres", connStr)
	check(err)

	err = db.Ping()
	check(err)

	AppLogger.Info("Database successfully connected")

	return db, nil
}
