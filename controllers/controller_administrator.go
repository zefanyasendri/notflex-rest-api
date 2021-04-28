package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

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
	generateToken(w, person.Email, person.Password, 0, 0)
	response.Status = 200
	response.Message = "Success Login <WELCOME>"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMemberBaseOnEmail(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()
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

	var response models.Response
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
	db := db.ConnectDB()

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
	db := database.ConnectDB()

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

	response := models.FilmResponse{Status: 200, Data: hasils, Message: "Data Found"}
	results, err := json.Marshal(response)
	

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}
