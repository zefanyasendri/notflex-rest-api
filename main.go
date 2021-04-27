package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/controllers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", controllers.LoginAdmin).Methods("GET")
	router.HandleFunc("/loginmember", controllers.Login).Methods("GET")
	router.HandleFunc("/watch/{id}", controllers.Authenticate(controllers.WatchFilm, 1)).Methods("GET")
	router.HandleFunc("/getuserbyemail", controllers.Authenticate(controllers.GetMemberBaseOnEmail, 0)).Methods("GET")
	router.HandleFunc("/suspend/{id}", controllers.SuspendMember).Methods("PUT")
	router.HandleFunc("/addfilm", controllers.AddFilm).Methods("POST")
	router.HandleFunc("/updatefilmbyid/{id}", controllers.UpdateFilmById).Methods("PUT")
	router.HandleFunc("/getfilmbykeyword/{keyword}", controllers.GetFilmByKeyword).Methods("GET")
	router.HandleFunc("/updateprofile/{id}", controllers.UpdateProfile).Methods("PUT")
	router.HandleFunc("/getfilmbyid/{id}", controllers.GetFilmByID).Methods("GET")

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
