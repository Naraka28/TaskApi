package main

import (
	"fmt"
	"go-server/database"
	"go-server/internal/middleware"
	"go-server/internal/task"
	"go-server/internal/user"
	"net/http"

	"github.com/rs/cors"
)
func main(){
	db, err := database.InitDb()
	if err != nil {
		fmt.Printf("Error initializing db connection: %v", err)
	}

	taskRepo := task.NewRepository(db)
	taskHandler := task.NewHandler(taskRepo)

	userRepo := user.NewRepository(db)
	userHandler := user.NewHandler(userRepo)



	mux := http.NewServeMux()

	mux.Handle("GET /tasks", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.GetAll)))
    mux.Handle("GET /tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.FindTaskById)))
    mux.Handle("POST /tasks", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.Save)))
    mux.Handle("DELETE /tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.Delete)))
    mux.Handle("DELETE /tasks", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.DeleteAll)))
    mux.Handle("PATCH /tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.Toggle)))
    mux.Handle("PUT /tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.Edit)))

	mux.Handle("GET /users",middleware.JWTMiddleware(http.HandlerFunc(userHandler.GetAll)))
    mux.HandleFunc("POST /users", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)


	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*", "http://127.0.0.1:5500"},
        AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH", "PUT"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

	handler := c.Handler(mux)

	fmt.Printf("Inicializando servidor en puerto %d\n",3000)
	http.ListenAndServe(":3000", handler)
}