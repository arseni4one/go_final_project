package database

import (
	"database/sql"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу и возвращает ID как строку
func AddTask(task Task) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// GetTask возвращает задачу по ID (строке)
func GetTask(idStr string) (Task, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return Task{}, err
	}
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

// UpdateTask обновляет задачу и возвращает ошибку, если запись не найдена
func UpdateTask(task Task) error {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteTask удаляет задачу по ID (строке)
func DeleteTask(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err = DB.Exec(query, id)
	return err
}

// UpdateTaskDate обновляет только дату (ID приходит как строка)
func UpdateTaskDate(idStr string, date string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = DB.Exec(query, date, id)
	return err
}

// Tasks возвращает список задач
func Tasks(limit int) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var id int
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		t.ID = strconv.Itoa(id)
		tasks = append(tasks, t)
	}

	// Проверяем ошибки после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
