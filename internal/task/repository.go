package task

import (
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (tr *TaskRepository) GetTasks(userID int) ([]Task, error) {
	rows, err := tr.db.Query("SELECT id, title, completed, userId FROM tasks WHERE userId = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener tareas: %v", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Completed, &task.UserId); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (tr *TaskRepository) Save(task TaskForm) (Task, error) {
	result, err := tr.db.Exec("INSERT INTO tasks(title, completed, userId) VALUES(?,?,?)", task.Title, task.Completed, task.UserId)
	if err != nil {
		return Task{}, fmt.Errorf("error al crear tarea: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf("error al obtener ID insertado: %v", err)
	}

	return tr.FindTaskById(int(id))
}

func (tr *TaskRepository) Delete(taskID int, userID int) error {
    result, err := tr.db.Exec("DELETE FROM tasks WHERE id = ? AND userId = ?", taskID, userID)
    if err != nil {
        return err
    }
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("tarea no encontrada o no tienes permiso")
    }
    return nil
}

func (tr *TaskRepository) ToggleTask(taskID int) (Task, error) {
    result, err := tr.db.Exec("UPDATE tasks SET completed = NOT completed WHERE id = ?", taskID)
    if err != nil {
        return Task{}, err
    }
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return Task{}, fmt.Errorf("permiso denegado o tarea inexistente")
    }
    return tr.FindTaskById(taskID)
}

func (tr *TaskRepository) DeleteAll(userID int) error {
	_, err := tr.db.Exec("DELETE FROM tasks WHERE userId = ?", userID)
	if err != nil {
		return fmt.Errorf("no se pudieron borrar las tareas: %v", err)
	}
	return nil
}

func (tr *TaskRepository) Edit(taskID int, title string) (Task, error) {
    result, err := tr.db.Exec("UPDATE tasks SET title = ? WHERE id = ?", title, taskID)
    if err != nil {
        return Task{}, err
    }
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return Task{}, fmt.Errorf("permiso denegado")
    }
    return tr.FindTaskById(taskID)
}

func (tr *TaskRepository) FindTaskById(id int) (Task, error) {
	var task Task
	err := tr.db.QueryRow("SELECT id, title, completed, userId FROM tasks WHERE id = ?", id).
		Scan(&task.Id, &task.Title, &task.Completed, &task.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("tarea no encontrada")
		}
		return task, err
	}
	return task, nil
}