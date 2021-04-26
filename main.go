package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/controllers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/loginasadminbroo", controllers.LoginAdmin).Methods("GET")
	router.HandleFunc("/login", controllers.LoginMember).Methods("GET")
	router.HandleFunc("/getuserbyemail", controllers.Authenticate(controllers.GetMemberBaseOnEmail, 0)).Methods("GET")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/logout", controllers.SignOut).Methods("GET")

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
