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

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	email := r.FormValue("email")
	password := r.FormValue("password")

	w.Header().Set("Content-Type", "application/json")
	if email == "admin" && password == "12345" {
		generateToken(w, email, password, 0, 0)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Status: http.StatusInternalServerError, Message: "Login Failed"})
		return
	}
	response = models.Response{
		Status:  http.StatusOK,
		Message: "Login success. Welcome " + email + "!",
		Data:    nil,
	}
	json.NewEncoder(w).Encode(response)
}
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
	var response models.Response
	if len(hasil.Email) == 0 {
		response = models.Response{Status: 404, Message: "Data Not Found"}
	} else {
		response = models.Response{Status: 200, Data: hasil, Message: "Data Found"}
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

	var response models.Response
	if err := db.Debug().Where("id_member = ?", idMember).Find(&member).Error; err != nil {
		response = models.Response{Status: 404, Message: "Member account not found"}
	} else {
		db.Model(&member).Where("status_akun LIKE ? OR status_akun LIKE ? AND id_member = ?", "%Subscribed%", "%Active%", idMember).Updates(memberUpdates)
		response = models.Response{Status: 200, Message: "Member account suspended"}
	}
	
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Menambah film baru ke database
func AddFilm(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	var film models.Film
	
	json.Unmarshal(body, &film)

	db.Debug().Create(&film)

	response := models.FilmResponse{Status: 200, Data: film, Message: "Added Film"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Mengubah data film sesuai id
func UpdateFilmById(w http.ResponseWriter, r *http.Request) {
	
	type result struct {
		IdFilm     int        `json:"idFilm"`
		Judul      string     `json:"judul"`
		TahunRilis string     `json:"tahunRilis"`
		Sutradara  string     `json:"sutradara"`
		Sinopsis   string     `json:"sinopsis"`
		IdGenre    int        `json:"idGenre"`
	}

	db := db.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	idFilm := vars["id"]

	var hasil result
	var filmUpdates models.Film
	
	json.Unmarshal(body, &filmUpdates)
	
	var film models.Film

	query := db.Table("films").Select("id_film, judul, tahun_rilis, sutradara, sinopsis, id_genre").Where("id_film = ?", idFilm).Row()
	db.Model(&film).Where("id_film = ?", idFilm).Updates(filmUpdates)

	query.Scan(&hasil.IdFilm, &hasil.Judul, &hasil.TahunRilis, &hasil.Sutradara, &hasil.Sinopsis, &hasil.IdGenre) 


	var response models.FilmResponse
	if hasil.IdFilm == 0 {
		response = models.FilmResponse{Status: 404, Message: "Film Dak Ditemukan"}
	} else {
		db.Table("films").Select("id_film, judul, tahun_rilis, sutradara, sinopsis, id_genre").Where("id_film = ?", idFilm).Find(&hasil)
		response = models.FilmResponse{Status: 200, Data: hasil, Message: "Film Data Updated"}
	}

	results, err := json.Marshal(response)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)

}

//Mengambil data film sesuai keyword yang diinputkan
func GetFilmByKeyword(w http.ResponseWriter, r *http.Request) {
	type result struct {
		IdFilm     int        `json:"idFilm"`
		Judul      string     `json:"judul"`
		TahunRilis string     `json:"tahunRilis"`
		Sutradara  string     `json:"sutradara"`
		Sinopsis   string     `json:"sinopsis"`
		IdGenre    int        `json:"idGenre"`
		JenisGenre string 	  `json:"JenisGenre"`
		NamaPemain []string   `json:"NamaPemain"`
	}
	db := db.ConnectDB()

	vars := mux.Vars(r)
	keywordJudul := vars["keyword"]
	
	var hasil result
	var hasils []result

	query_film, _ := db.Debug().Table("films").Select("films.id_film, films.judul, films.tahun_rilis, films.sutradara, films.sinopsis, films.id_genre, genres.jenis_genre").Joins("JOIN genres ON films.id_genre = genres.id_genre").Where("films.judul LIKE ?", "%"+keywordJudul+"%").Rows()

	defer query_film.Close()

	for query_film.Next() {

		query_film.Scan(&hasil.IdFilm, &hasil.Judul, &hasil.TahunRilis, &hasil.Sutradara, &hasil.Sinopsis, &hasil.IdGenre, &hasil.JenisGenre)

		query_pemain, _ := db.Debug().Table("pemains").Select("pemains.nama_pemain, list_pemains.peran").Joins("JOIN list_pemains ON pemains.id_pemain = list_pemains.id_pemain").Joins("JOIN films ON list_pemains.id_film = films.id_film").Where("films.id_film = ?", &hasil.IdFilm).Rows()
	
		for query_pemain.Next() {
			var pemain string
			var peranPemain string
			query_pemain.Scan(&pemain, &peranPemain)
			if pemain != "" {
				hasil.NamaPemain = append(hasil.NamaPemain, pemain)
				hasil.NamaPemain = append(hasil.NamaPemain, peranPemain)
			}
		}
		hasils = append(hasils, hasil)
		hasil.NamaPemain = nil
	}
	//var results []byte
	var response models.FilmResponse
	//var err error
	if len(hasil.Judul) == 0 {
		response = models.FilmResponse{Status: 400, Message: "Data Not Found"}
	} else {
		response = models.FilmResponse{Status: 200, Data: hasils, Message: "Data Found"}
	}

	results, err := json.Marshal(response)	

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

//Mengambil data film berdasarkan ID
func GetFilmByID(w http.ResponseWriter, r *http.Request) {
	type result struct {
		IdFilm     int      `json:"idFilm"`
		Judul      string   `json:"judul"`
		TahunRilis string   `json:"tahunRilis"`
		Sutradara  string   `json:"sutradara"`
		Sinopsis   string   `json:"sinopsis"`
		IdGenre    int      `json:"idGenre"`
		JenisGenre string   `json:"JenisGenre"`
		NamaPemain []string `json:"NamaPemain"`
	}
	db := db.ConnectDB()

	vars := mux.Vars(r)
	id := vars["id"]

	var hasil result
	var hasils []result

	query_film, _ := db.Debug().Table("films").Select("films.id_film, films.judul, films.tahun_rilis, films.sutradara, films.sinopsis, films.id_genre, genres.jenis_genre").Joins("JOIN genres ON films.id_genre = genres.id_genre").Where("films.id_film = ?", id).Rows()

	defer query_film.Close()

	for query_film.Next() {

		query_film.Scan(&hasil.IdFilm, &hasil.Judul, &hasil.TahunRilis, &hasil.Sutradara, &hasil.Sinopsis, &hasil.IdGenre, &hasil.JenisGenre)

		query_pemain, _ := db.Debug().Table("pemains").Select("pemains.nama_pemain, list_pemains.peran").Joins("JOIN list_pemains ON pemains.id_pemain = list_pemains.id_pemain").Joins("JOIN films ON list_pemains.id_film = films.id_film").Where("films.id_film = ?", &hasil.IdFilm).Rows()

		for query_pemain.Next() {
			var pemain string
			var peranPemain string
			query_pemain.Scan(&pemain, &peranPemain)
			if pemain != "" {
				hasil.NamaPemain = append(hasil.NamaPemain, pemain)
				hasil.NamaPemain = append(hasil.NamaPemain, peranPemain)
			}
		}
		hasils = append(hasils, hasil)
		hasil.NamaPemain = nil
	}

	response := models.FilmResponse{Status: 200, Data: hasils, Message: "Data Found"}
	results, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}
