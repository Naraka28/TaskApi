package task

import (
	"encoding/json"
	"fmt"
	"go-server/utils"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	repo *TaskRepository
}

func NewHandler(tr *TaskRepository) *TaskHandler {
	return &TaskHandler{repo: tr}
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Extraemos el ID del contexto (puesto por tu middleware de Auth)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendJSONError(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	tasks, err := h.repo.GetTasks(userID)
	if err != nil {
		utils.SendJSONError(w, "Error al buscar tareas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Save(w http.ResponseWriter, r *http.Request) {
	var newTask TaskForm
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		utils.SendJSONError(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Forzamos que el UserId de la tarea sea el del usuario logueado
	if userID, ok := r.Context().Value("user_id").(int); ok {
		newTask.UserId = userID
	}

	task, err := h.repo.Save(newTask)
	if err != nil {
		utils.SendJSONError(w, "No se pudo guardar la tarea", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.SendJSONError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	task, err := h.repo.ToggleTask(id)
	if err != nil {
		utils.SendJSONError(w, "No se pudo cambiar el estado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var data UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.SendJSONError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	task, err := h.repo.Edit(id,data.Title)
	if err != nil {
		utils.SendJSONError(w, "Error al editar tarea", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
    taskID, _ := strconv.Atoi(r.PathValue("id"))
    userID, _ := r.Context().Value("user_id").(int)

    if err := h.repo.Delete(taskID, userID); err != nil {
        utils.SendJSONError(w, err.Error(), http.StatusForbidden)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{"status": "deleted"})
}

func (h *TaskHandler) FindTaskById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.SendJSONError(w, "ID debe ser un número", http.StatusBadRequest)
		return
	}

	task, err := h.repo.FindTaskById(id)
	if err != nil {
		utils.SendJSONError(w, fmt.Sprintf("Tarea %d no encontrada", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)

	if err := h.repo.DeleteAll(userID); err != nil {
		utils.SendJSONError(w, "No se pudieron borrar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Todas las tareas borradas"})
}