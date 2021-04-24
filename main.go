package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/controllers"
	db "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
)

func main() {
	db.ConnectDB()
	router := mux.NewRouter()

	router.HandleFunc("/login", controllers.LoginAdmin).Methods("GET")
	router.HandleFunc("/getuserbyemail", controllers.Authenticate(controllers.GetMemberBaseOnEmail, 0)).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 4321")
	log.Fatal(http.ListenAndServe(":4321", router))
}
