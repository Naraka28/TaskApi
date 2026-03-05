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

type Response struct{
    Data []Task
}

var tasks = []Task{
        {1,"Lavarme los dientes",false},
        {2,"Terminar tarea de Frontend",false},
        {3,"Preparar unos taquitos de carne asada aca bien sabrosos",true},
        {4,"Ir al gimnasio",false},
        {5,"Tomar un bañito y echarme una pestañita",true},
    }

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", getTasks)
    mux.HandleFunc("GET /tasks/{id}", getTaskByID)
    mux.HandleFunc("POST /tasks", createTask)
    mux.HandleFunc("DELETE /tasks/{id}", deleteTask)

	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://127.0.0.1:5500"},
        AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
        AllowCredentials: true,
    })

	handler := c.Handler(mux)

	fmt.Printf("Inicializando servidor en puerto %d",80)
	http.ListenAndServe(":80", handler)
}
func getTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)

}
func createTask(w http.ResponseWriter, r *http.Request){
	var newTask Task

	error := json.NewDecoder(r.Body).Decode(&newTask)

	if error != nil{
		message := error.Error()
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTask)


}
func deleteTask(w http.ResponseWriter, r *http.Request){
	param := r.PathValue("id")
	id, error := strconv.Atoi(param)
	if error != nil {
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
	tasks = append(tasks[:i], tasks[i+1:]...)

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

func getTaskByID(w http.ResponseWriter, r *http.Request){
	value := r.PathValue("id")
	id, error := strconv.Atoi(value)
	if error != nil{
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