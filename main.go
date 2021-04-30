package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/controllers"
)

func main() {
	router := mux.NewRouter()

	// Hans
	router.HandleFunc("/loginadmin", controllers.LoginAdmin).Methods("GET")
	router.HandleFunc("/regis", controllers.Register).Methods("POST")
	router.HandleFunc("/getuserbyemail", controllers.Authenticate(controllers.GetMemberBaseOnEmail, 0)).Methods("GET")
	router.HandleFunc("/logout", controllers.SignOut).Methods("GET")

	// Nealson
	router.HandleFunc("/suspend/{id}", controllers.Authenticate(controllers.SuspendMember, 0)).Methods("PUT")
	router.HandleFunc("/addfilm", controllers.Authenticate(controllers.AddFilm, 0)).Methods("POST")
	router.HandleFunc("/updatefilmbyid/{id}", controllers.Authenticate(controllers.UpdateFilmById, 0)).Methods("PUT")
	router.HandleFunc("/getfilmbykeyword/{keyword}", controllers.Authenticate(controllers.GetFilmByKeyword, 0)).Methods("GET")

	// Zefanya
	router.HandleFunc("/updateprofile", controllers.Authenticate(controllers.UpdateProfile, 1)).Methods("PUT")
	router.HandleFunc("/getfilmbyid/{id}", controllers.Authenticate(controllers.GetFilmByID, 0)).Methods("GET")
	router.HandleFunc("/getfilmbykeywords/{keywords}", controllers.Authenticate(controllers.GetFilmByKeywords, 1)).Methods("GET")
	router.HandleFunc("/getwatchhistory", controllers.Authenticate(controllers.GetWatchHistory, 1)).Methods("GET")

	//Hilbert
	router.HandleFunc("/loginmember", controllers.Login).Methods("GET")
	router.HandleFunc("/watch/{id}", controllers.Authenticate(controllers.WatchFilm, 1)).Methods("GET")
	router.HandleFunc("/subscribe", controllers.Authenticate(controllers.Subscribe, 1)).Methods("PUT")
	router.HandleFunc("/unsubscribe", controllers.Authenticate(controllers.Unsubscribe, 1)).Methods("PUT")

	fmt.Println("Connected to port 4321")
	log.Fatal(http.ListenAndServe(":4321", router))
}
