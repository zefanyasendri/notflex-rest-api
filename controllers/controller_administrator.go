package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	database "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	pass := r.URL.Query()["password"]
	email := r.URL.Query()["email"]
	var person models.Person
	var response models.PersonResponse

	if err := db.Where("email = ? and password = ?", email[0], pass[0]).First(&person).Error; err != nil {
		log.Print(err)
		response.Status = 400
		response.Message = "Error"
		return
	}

	db.Find(&person)
	generateToken(w, person.Email, person.Password, 0)
	response.Status = 200
	response.Message = "Success Login <WELCOME>"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMemberBaseOnEmail(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()

	email := r.URL.Query()["email"]

	query := "SELECT * FROM member"

	if email != nil {
		query += " WHERE email ='" + email[0] + "'"
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Print(err)
	}

	var member models.Member
	var members []models.Member
	for rows.Next() {
		if err := rows.Scan(&member.IdMember, &member.NamaLengkap, &member.TanggalLahir, &member.JenisKelamin, &member.AsalNegara, &member.StatusAkun, &member.NoKartuKredit, &member.Email); err != nil {
			log.Print(err.Error())
		} else {
			members = append(members, member)
		}
	}

	var response models.MemberResponse
	if err == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = members
	} else {
		response.Status = 400
		response.Message = "Error"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SuspendMember(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	idMember := vars["id"]

	var memberUpdates models.Member
	json.Unmarshal(body, &memberUpdates)

	var member models.Member
	db.Where("WHERE status_akun = ? AND id_member = ?", "Active", idMember).Find(&member)
	db.Model(&member).Updates(memberUpdates)

	response := models.FilmResponse{Status: 200, Data: member, Message: "Member account suspended"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db.Save(&member)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func AddFilm(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	var film models.Film
	json.Unmarshal(body, &film)

	db.Create(&film)

	response := models.FilmResponse{Status: 200, Data: film, Message: "Added Film"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db.Save(&film)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UpdateFilmById(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	idFilm := vars["id"]

	var filmUpdates models.Film
	json.Unmarshal(body, &filmUpdates)

	var film models.Film
	db.Find(&film, idFilm)
	db.Model(&film).Updates(filmUpdates)

	response := models.FilmResponse{Status: 200, Data: film, Message: "Film Data Updated"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db.Save(&film)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func GetFilmByKeyword(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	vars := mux.Vars(r)
	keyword := vars["keyword"]

	var film []models.Film
	db.Where("judul LIKE ?", "%"+keyword+"%").Find(&film)

	response := models.FilmResponse{Status: 200, Data: film, Message: "Data Found"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetFilmByID(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	vars := mux.Vars(r)
	id := vars["id"]

	var film []models.Film
	// db.Where("judul = ?", id).Find(&film)
	db.First(&film, id)

	response := models.FilmResponse{Status: 200, Data: film, Message: "Data Found"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
