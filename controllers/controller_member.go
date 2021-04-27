package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	database "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	vars := mux.Vars(r)
	id_member := vars["id"]

	body, _ := ioutil.ReadAll(r.Body)

	var profileUpdates models.Member
	json.Unmarshal(body, &profileUpdates)

	var member models.Member
	db.Find(&member, id_member)
	db.Model(&member).Updates(profileUpdates)

	response := models.MemberResponse{Status: 200, Data: member, Message: "Member Data Updated"}
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
	db := database.ConnectDB()

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
