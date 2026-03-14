package task

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskHandler struct{
	repo *TaskRepository
}
func NewHandler(tr *TaskRepository) *TaskHandler{
	return &TaskHandler{repo: tr}
}

func (handler *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("user_id").(int)

	tasks, err := handler.repo.GetTasks(userID)
	if err != nil {
		fmt.Printf("Error Handling Tasks: %v", err)
		return
	}

	w.Header().Set("Content-Type:", "application/json")
	json.NewEncoder(w).Encode(tasks)
}