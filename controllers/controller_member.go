package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	email := r.FormValue("email")
	password := r.FormValue("password")

	w.Header().Set("Content-Type", "application/json")

	err := CheckLogin(email, password)
	if err != nil {
		response = models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Email/password seems to be incorrect. Please try again.",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	generateToken(w, email, password, 1)

	res2, err := CheckSuspended(email)

	if err != nil {
		response = models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again.",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !res2 {
		response = models.Response{
			Status:  http.StatusUnauthorized,
			Message: "I'm sorry looks like your account is been suspended.",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response = models.Response{
		Status:  http.StatusOK,
		Message: "Login success. Welcome " + email + "!",
		Data:    nil,
	}
	json.NewEncoder(w).Encode(response)
}

func CheckSuspended(email string) (bool, error) {
	var member models.Member

	db := db.ConnectDB()
	row := db.Table("members").Where("email = ?", email).Select("status_akun").Row()
	err := row.Err()

	if err != nil {
		fmt.Println("Query error")
		return false, err
	}

	row.Scan(&member.StatusAkun)
	status := member.StatusAkun

	if status == "Suspended" {
		return false, nil
	}

	return true, nil
}

func CheckLogin(email, password string) error {
	var pwd string
	db := db.ConnectDB()
	row := db.Table("members").Where("email = ?", email).Select("password").Row()
	err := row.Err()

	if err != nil {
		fmt.Println("Query error")
		return err
	}

	row.Scan(&pwd)
	match, err := CheckHashedPassword(password, pwd)
	if !match {
		fmt.Println("Hash and password does not match.")
		return err
	}
	return nil
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

	vars := mux.Vars(r)
	id_member := vars["id"]

	body, _ := ioutil.ReadAll(r.Body)

	var profileUpdates models.Member
	json.Unmarshal(body, &profileUpdates)

	var member models.Member
	db.Find(&member, id_member)
	db.Model(&member).Updates(profileUpdates)

	response := models.Response{Status: 200, Data: member, Message: "Member Data Updated"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db.Save(&member)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetFilmByKeywords(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

	vars := mux.Vars(r)
	keywordfilm := vars["keywords"]

	var film []models.Film

	db.Table("films").Select("films.id_film,films.judul,films.tahun_rilis, films.sutradara, films.sinopsis, films.id_genre").Joins("LEFT JOIN genres ON films.id_genre = genres.id_genre LEFT JOIN list_pemains ON films.id_film = list_pemains.id_film LEFT JOIN pemains ON list_pemains.id_pemain = pemains.id_pemain").Where("films.judul LIKE ? OR films.sutradara LIKE ? OR films.tahun_rilis LIKE ? OR films.sinopsis LIKE ? OR genres.jenis_genre LIKE ? OR pemains.nama_pemain LIKE ?", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%").Scan(&film)

	response := models.FilmResponse{Status: 200, Data: film, Message: "Data Found"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetWatchHistory(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Judul string    `json:"judul"`
		Waktu time.Time `json:"waktu"`
	}

	db := db.ConnectDB()

	vars := mux.Vars(r)
	id_member := vars["id"]

	query, err := db.Table("films").Select("films.judul, histories.tanggal_nonton").Joins("LEFT JOIN histories ON histories.id_film = films.id_film LEFT JOIN members ON histories.id_member = members.id_member").Where("histories.id_member = ?", id_member).Rows()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer query.Close()
	var response Response
	var responses []Response

	for query.Next() {
		query.Scan(&response.Judul, &response.Waktu)
		responses = append(responses, response)
	}

	res := models.FilmResponse{Status: 200, Data: responses, Message: "Data Found"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
