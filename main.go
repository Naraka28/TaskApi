package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)
type Task struct{
    Id int `json:"id"`
    Title string `json:"title"`
    Completed bool `json:"completed"`
}

type UpdateTitle struct {
	Title string `json:"title"`
}


var tasks = []Task{
        {1,"Lavarme los dientes",false},
        {2,"Terminar tarea de Frontend",false},
        {3,"Preparar unos taquitos de carne asada aca bien sabrosos",false},
        {4,"Ir al gimnasio",false},
        {5,"Tomar un bañito y echarme una pestañita",true},
    }

var completedTask []Task
var contador = 5

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", getTasks)
    mux.HandleFunc("GET /tasks/{id}", getTaskByID)
    mux.HandleFunc("POST /tasks", createTask)
    mux.HandleFunc("DELETE /tasks/{id}", deleteTask)
    mux.HandleFunc("DELETE /tasks", deleteAll)
    mux.HandleFunc("PATCH /tasks/{id}", toggleTask)
    mux.HandleFunc("PUT /tasks/{id}", editTask)

	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://127.0.0.1:5500"},
        AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH", "PUT"},
        AllowedHeaders:   []string{"Content-Type"},
        AllowCredentials: true,
    })

	handler := c.Handler(mux)

	fmt.Printf("Inicializando servidor en puerto %d\n",3000)
	http.ListenAndServe(":3000", handler)
}
func getTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)

}
func createTask(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)

	if err != nil{
		message := err.Error()
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	newTask.Id = contador + 1
	contador++
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTask)


}
func deleteTask(w http.ResponseWriter, r *http.Request){
	param := r.PathValue("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		message := fmt.Sprintf("Couldnt convert %s", param)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	i, ok := search(id, tasks)

	if !ok{
		message := fmt.Sprintf("Coulndt delete taks with ID: %d",id)
		http.Error(w, message, http.StatusNotFound)
		return
	}
	deleteTodo(i, &tasks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}
func search(id int, t []Task) (index int, ok bool){
    for i, task := range t{
        if task.Id == id {
            return i, true
        }
    }
    return -1, false
}
func deleteTodo(index int, t *[]Task){
	*t = append((*t)[:index], (*t)[index+1:]...)

}

func getTaskByID(w http.ResponseWriter, r *http.Request){
	value := r.PathValue("id")
	id, err := strconv.Atoi(value)
	if err != nil{
		http.Error(w, "Not a number", http.StatusBadRequest)
		return
	}

	i, ok := search(id, tasks);

	if !ok{
		mensaje := fmt.Sprintf("Task with ID: %d not found", id)
		http.Error(w, mensaje, http.StatusNotFound)
		return
	}

	task := tasks[i]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

func toggleTask(w http.ResponseWriter, r *http.Request){
	value := r.PathValue("id")
	id, err := strconv.Atoi(value)
	if err != nil{
		http.Error(w, "Not a number", http.StatusBadRequest)
		return
	}

	i, ok := search(id, tasks);

	if !ok{
		mensaje := fmt.Sprintf("Task with ID: %d not found", id)
		http.Error(w, mensaje, http.StatusNotFound)
		return
	}
	tasks[i].Completed = !tasks[i].Completed

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func deleteAll(w http.ResponseWriter, r *http.Request){
	tasks = nil

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func editTask(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	var body UpdateTitle
	value := r.PathValue("id")
	err := json.NewDecoder(r.Body).Decode(&body)

	id, err := strconv.Atoi(value)
	if err != nil{
		http.Error(w, "Not a number", http.StatusBadRequest)
		return
	}

	i, ok := search(id, tasks);

	if !ok{
		mensaje := fmt.Sprintf("Task with ID: %d not found", id)
		http.Error(w, mensaje, http.StatusNotFound)
		return
	}
	tasks[i].Title = body.Title
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks[i])	

}