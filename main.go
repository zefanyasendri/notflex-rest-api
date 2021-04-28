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
	router.HandleFunc("/regis", controllers.Register).Methods("POST")
	router.HandleFunc("/getuserbyemail", controllers.GetMemberBaseOnEmail).Methods("GET")

	// Nealson
	router.HandleFunc("/suspend/{id}", controllers.SuspendMember).Methods("PUT")
	router.HandleFunc("/addfilm", controllers.AddFilm).Methods("POST")
	router.HandleFunc("/updatefilmbyid/{id}", controllers.UpdateFilmById).Methods("PUT")
	router.HandleFunc("/getfilmbykeyword/{keyword}", controllers.GetFilmByKeyword).Methods("GET")

	// Zefa
	router.HandleFunc("/updateprofile/{id}", controllers.UpdateProfile).Methods("PUT")
	router.HandleFunc("/getfilmbyid/{id}", controllers.GetFilmByID).Methods("GET")
	router.HandleFunc("/getfilmbykeywords/{keywords}", controllers.GetFilmByKeywords).Methods("GET")
	router.HandleFunc("/getwatchhistory", controllers.Authenticate(controllers.GetWatchHistory, 1)).Methods("GET")

	//Hilbert
	router.HandleFunc("/loginmember", controllers.Login).Methods("GET")
	router.HandleFunc("/watch/{id}", controllers.Authenticate(controllers.WatchFilm, 1)).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(router)

	fmt.Println(controllers.HashPassword("12345"))

	http.Handle("/", handler)
	fmt.Println("Connected to port 4321")
	log.Fatal(http.ListenAndServe(":4321", router))
}
