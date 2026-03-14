package task

import (
	"database/sql"
	"fmt"
)

type TaskRepository struct{
	r *sql.DB
}

func NewRepository(db *sql.DB) *TaskRepository{
	return &TaskRepository{
		r: db,
	}
}

func (tr *TaskRepository) GetTasks(id int) ([]Task, error){
	var tasks []Task
	rows, err := tr.r.Query("SELECT tasks.id, title, completed, userId FROM tasks INNER JOIN users ON tasks.userId = ?;", id)

	if err != nil {
		return nil, fmt.Errorf("Retrieving Tasks: %v", err)
	}

	defer rows.Close()

	for rows.Next(){
		var task Task
		if err = rows.Scan(&task.Id, &task.Title, &task.Completed, &task.UserId); err != nil {
			return nil, fmt.Errorf("Retrieving Tasks: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Retrieving Tasks: %v", err)
	}
	return tasks, nil

}
// func createTask(w http.ResponseWriter, r *http.Request){
// 	defer r.Body.Close()

// 	var newTask Task

// 	err := json.NewDecoder(r.Body).Decode(&newTask)

// 	if err != nil{
// 		message := err.Error()
// 		http.Error(w, message, http.StatusBadRequest)
// 		return
// 	}
// 	newTask.Id = contador + 1
// 	contador++
// 	tasks = append(tasks, newTask)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)

// 	json.NewEncoder(w).Encode(newTask)


// }
// func deleteTask(w http.ResponseWriter, r *http.Request){
// 	param := r.PathValue("id")
// 	id, err := strconv.Atoi(param)
// 	if err != nil {
// 		message := fmt.Sprintf("Couldnt convert %s", param)
// 		http.Error(w, message, http.StatusBadRequest)
// 		return
// 	}

// 	i, ok := search(id, tasks)

// 	if !ok{
// 		message := fmt.Sprintf("Coulndt delete taks with ID: %d",id)
// 		http.Error(w, message, http.StatusNotFound)
// 		return
// 	}
// 	deleteTodo(i, &tasks)

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(tasks)

// }
// func search(id int, t []Task) (index int, ok bool){
//     for i, task := range t{
//         if task.Id == id {
//             return i, true
//         }
//     }
//     return -1, false
// }
// func deleteTodo(index int, t *[]Task){
// 	*t = append((*t)[:index], (*t)[index+1:]...)

// }

// func getTaskByID(w http.ResponseWriter, r *http.Request){
// 	value := r.PathValue("id")
// 	id, err := strconv.Atoi(value)
// 	if err != nil{
// 		http.Error(w, "Not a number", http.StatusBadRequest)
// 		return
// 	}

// 	i, ok := search(id, tasks);

// 	if !ok{
// 		mensaje := fmt.Sprintf("Task with ID: %d not found", id)
// 		http.Error(w, mensaje, http.StatusNotFound)
// 		return
// 	}

// 	task := tasks[i]
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(task)

// }

// func toggleTask(w http.ResponseWriter, r *http.Request){
// 	value := r.PathValue("id")
// 	id, err := strconv.Atoi(value)
// 	if err != nil{
// 		http.Error(w, "Not a number", http.StatusBadRequest)
// 		return
// 	}

// 	i, ok := search(id, tasks);

// 	if !ok{
// 		mensaje := fmt.Sprintf("Task with ID: %d not found", id)
// 		http.Error(w, mensaje, http.StatusNotFound)
// 		return
// 	}
// 	tasks[i].Completed = !tasks[i].Completed

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(tasks)
// }

// func deleteAll(w http.ResponseWriter, r *http.Request){
// 	tasks = nil

// 	w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(tasks)
// }

// func editTask(w http.ResponseWriter, r *http.Request){
// 	defer r.Body.Close()
// 	var body UpdateTitle
// 	value := r.PathValue("id")
// 	err := json.NewDecoder(r.Body).Decode(&body)

// 	id, err := strconv.Atoi(value)
// 	if err != nil{
// 		http.Error(w, "Not a number", http.StatusBadRequest)
// 		return
// 	}

// 	i, ok := search(id, tasks);

// 	if !ok{
// 		mensaje := fmt.Sprintf("Task with ID: %d not found", id)
// 		http.Error(w, mensaje, http.StatusNotFound)
// 		return
// 	}
// 	tasks[i].Title = body.Title
// 	w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(tasks[i])	

// }