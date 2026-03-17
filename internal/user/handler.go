package user

import (
	"encoding/json"
	"fmt"
	"go-server/internal/auth"
	"go-server/utils"
	"net/http"
)

type UserHandler struct{
	repo *UserRepository
}

func NewHandler(ur *UserRepository) *UserHandler{
	return &UserHandler{ repo : ur}
}

func (handler *UserHandler) GetAll(w http.ResponseWriter, r * http.Request){
	users, err := handler.repo.GetAll()
	if err != nil {

	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)

}

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	var newUser RegisterUser

	err := json.NewDecoder(r.Body).Decode(&newUser)

	authRepo := auth.NewRepository(handler.repo.r)
	newUser.Password, err = authRepo.HashPassword(newUser.Password)

	if err != nil {
		fmt.Print("Error hashing")
		return
	}

	id, err := handler.repo.Register(newUser)

	if err != nil{
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request){
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
	authRepo := auth.NewRepository(handler.repo.r)

    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        utils.SendJSONError(w, "Datos inválidos", http.StatusBadRequest)
        return
    }

    u, err := handler.repo.GetByEmail(credentials.Email)
    if err != nil {
        utils.SendJSONError(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
        return
    }

    if !authRepo.CheckPasswordHash(credentials.Password, u.Password) {
        utils.SendJSONError(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
        return
    }

    token, err := authRepo.GenerateJWT(fmt.Sprintf("%d", u.Id), "esternocleidomastoideo")
    if err != nil {
        utils.SendJSONError(w, "Error al generar el token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}