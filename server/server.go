package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MirekKrassilnikov/Todo_list_api_sql/repeater"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type Errr struct {
	error string `json:"error,omitempty"`
}

type Controller struct {
	DB *sql.DB
}

func (ctl *Controller) TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		ctl.HandlePost(w, r) // Обработка POST-запросов // Обработка POST-запросов

	case http.MethodGet:
		ctl.getTaskById(w, r) // Обработка GET-запросов
	case http.MethodPut:
		ctl.UpdateTask(w, r)

	case http.MethodDelete:
		ctl.DeleteTaskByID(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (ctl *Controller) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	codeAndNumber := strings.Split(task.Repeat, " ")

	// Проверка обязательного поля id и title
	if len(task.ID) == 0 {
		http.Error(w, `{"error":"ID is required"}`, http.StatusBadRequest)
		return
	}
	_, err = strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error":"ошибка конвертации"}`, http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		return
	}
	if len(codeAndNumber) == 0 || (codeAndNumber[0] != "y" && codeAndNumber[0] != "d") {
		http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	_, err = time.Parse(repeater.Layout, task.Date)
	if err != nil {
		http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли задача
	var existingTask Task
	err = ctl.DB.QueryRow("SELECT id FROM scheduler WHERE id = ?", task.ID).Scan(&existingTask.ID)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
		return
	}

	// Выполняем обновление задачи
	updateSQL := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = ctl.DB.Exec(updateSQL, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error":"Failed to update task"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем пустой JSON в случае успешного обновления
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func (ctl *Controller) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := ctl.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		http.Error(w, `{"error":"error with rows"}`, http.StatusNotFound)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"error with rows"}`, http.StatusNotFound)
			return
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = []Task{}
	}
	responseMap := map[string][]Task{"tasks": tasks}
	response, err := json.Marshal(responseMap)
	if err != nil {
		http.Error(w, "Response JSON creation error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (ctl *Controller) getTaskById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		http.Error(w, `{"error": "No ID provided"}`, http.StatusBadRequest)
		return
	}

	// Выполняем SQL-запрос для получения задачи по id
	var task Task
	err := ctl.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", idTask).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "Server error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Если задача найдена, возвращаем её в виде JSON
	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, `{"error": "Response JSON creation error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (ctl *Controller) HandlePost(w http.ResponseWriter, r *http.Request) {
	var task Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Проверка обязательного поля title
	if len(task.Title) == 0 {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты и установка текущей даты, если дата некорректна
	if task.Date == "" {
		task.Date = time.Now().Format(repeater.Layout)
	}
	timeTimeDate, err := time.Parse(repeater.Layout, task.Date)
	if err != nil {
		http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
		return
	}
	timeNowDateOnly := time.Now().Truncate(24 * time.Hour)
	if timeTimeDate.Before(timeNowDateOnly) {
		// Если дата задачи меньше текущей даты и есть правило повторения
		if task.Repeat != "" {
			nextDate, err := repeater.NextDate(time.Now().Format(repeater.Layout), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Invalid repeat rule"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(repeater.Layout)
		}
	}

	insertSQL := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?);`
	result, err := ctl.DB.Exec(insertSQL, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Failed to insert task"}`, http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Failed to retrieve task ID"}`, http.StatusInternalServerError)
		return
	}
	response := Response{
		ID: int(id),
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error":"Failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func MainHandle(res http.ResponseWriter, req *http.Request) {
	out := "Hello from server package!"
	res.Write([]byte(out))
}
func (ctl *Controller) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	var task Task
	err := ctl.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"server error"}`, http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		err = ctl.DeleteTaskByID(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		//w.Write([]byte(`{}`))
		return
	}
	now := time.Now().Format("20060102")
	nextDate, err := repeater.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"error calculating next date"}`, http.StatusInternalServerError)
		return
	}

	updateSQL := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = ctl.DB.Exec(updateSQL, nextDate, task.ID)

	if err != nil {
		http.Error(w, `{"error":"Failed to update task"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем пустой JSON в случае успешного обновления
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))

}

func (ctl *Controller) DeleteTaskByID(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"error calculating next date"}`, http.StatusInternalServerError)
		return nil
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"invalid id format"}`, http.StatusInternalServerError)
		return nil
	}
	_, err = ctl.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, `{"error":"failed to delete task"}`, http.StatusInternalServerError)
	}

	// Если задача успешно удалена, возвращаем пустой JSON {}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
	return nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"error": message,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}
func (ctl *Controller) ApiNextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Вызываем функцию NextDate
	nextDate, err := repeater.NextDate(nowStr, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
