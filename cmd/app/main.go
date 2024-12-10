package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"hakaton/internal/handlers"
	//"hakaton/internal/middleware"
)

func main() {
	r := mux.NewRouter()

	// Регистрация маршрутов
	r.HandleFunc("/api/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginUser).Methods("POST")

	// Маршруты с защитой (middleware JWT)
	api := r.PathPrefix("/api").Subrouter()
	//api.Use(middleware.JWTMiddleware)
	api.HandleFunc("/companies", handlers.GetCompanies).Methods("GET")
	api.HandleFunc("/companies/{id}", handlers.GetCompanyByID).Methods("GET")

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
