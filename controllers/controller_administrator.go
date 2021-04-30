package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func GetMemberBaseOnEmail(w http.ResponseWriter, r *http.Request) {

	type result struct {
		Email         string   `json:"email"`
		Password      string   `json:"password"`
		IdMember      int      `json:"idMember" gorm:"primaryKey"`
		NamaLengkap   string   `json:"namaLengkap"`
		TanggalLahir  string   `json:"tanggalLahir"`
		JenisKelamin  string   `json:"jenisKelamin"`
		AsalNegara    string   `json:"asalNegara"`
		StatusAkun    string   `json:"statusAkun"`
		NoKartuKredit string   `json:"noKartuKredit" gorm:"type:varchar(191)"`
		History       []string `json:"history"`
	}

	db := db.ConnectDB()

	email := r.URL.Query()["email"]

	var hasil result

	query_member, _ := db.Debug().Table("members").Select("*").Where("email = ?", email[0]).Rows()

	for query_member.Next() {
		query_member.Scan(&hasil.Email, &hasil.Password, &hasil.IdMember, &hasil.NamaLengkap, &hasil.TanggalLahir, &hasil.JenisKelamin, &hasil.AsalNegara, &hasil.StatusAkun, &hasil.NoKartuKredit)
		query_history, _ := db.Debug().Table("films").Select("films.judul, histories.tanggal_nonton").Joins("JOIN histories ON films.id_film = histories.id_film JOIN members ON histories.id_member = members.id_member").Where("members.email = ?", email[0]).Rows()
		for query_history.Next() {
			var Judulfilm string
			var TanggalNontonFilm string
			query_history.Scan(&Judulfilm, &TanggalNontonFilm)
			hasil.History = append(hasil.History, Judulfilm)
			hasil.History = append(hasil.History, TanggalNontonFilm)
			fmt.Println(hasil.History)
		}
	}
	var response models.MemberResponse
	if len(hasil.Email) == 0 {
		response = models.MemberResponse{Status: 404, Message: "Data Not Found"}
	} else {
		response = models.MemberResponse{Status: 200, Data: hasil, Message: "Data Found"}
	}
	results, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

//Suspend/Block status keaktifan member
func SuspendMember(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	idMember := vars["id"]

	var memberUpdates models.Member
	json.Unmarshal(body, &memberUpdates)

	var member models.Member
	db.Where("status_akun = ? AND id_member = ?", "Active", idMember).Find(&member)
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

//Menambah film baru ke database
func AddFilm(w http.ResponseWriter, r *http.Request) {
	//db := db.ConnectDB()
	type Hasilinput struct {
		IdFilm     int    `json:"idFilm"`
		Judul      string `json:"judul"`
		TahunRilis string `json:"tahunRilis"`
		Sutradara  string `json:"sutradara"`
		Sinopsis   string `json:"sinopsis"`
		IdGenre    int    `json:"idGenre"`
		IdPemain   int    `json:"idPemain"`
		NamaPemain string `json:"namaPemain"`
		Peran      string `json:"peran"`
	}
	body, _ := ioutil.ReadAll(r.Body)

	var input Hasilinput
	json.Unmarshal(body, &input)

	//db.Create(&film)
	fmt.Println(input)

	response := models.FilmResponse{Status: 200, Data: input, Message: "Added Film"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//db.Save(&film)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Mengubah data film sesuai id
func UpdateFilmById(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

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
