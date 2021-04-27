package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	email := r.FormValue("email")
	password := r.FormValue("password")

	w.Header().Set("Content-Type", "application/json")

	res1, err := CheckLogin(email, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status: http.StatusInternalServerError,
			// Message: "Email/password seems to be incorrect. Please try again.",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	res2, err := CheckSuspended(email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again.",
			Data:    nil,
		})
		return
	}

	if !res2 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusUnauthorized,
			Message: "I'm sorry looks like your account is been suspended.",
			Data:    nil,
		})
		return
	}

	generateToken(w, email, password, res1, 1)

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

func CheckLogin(email, password string) (int, error) {
	var pwd string
	var id_member int
	db := db.ConnectDB()
	row := db.Debug().Table("members").Where("email = ?", email).Select("id_member", "password").Row()
	err := row.Err()

	if err != nil {
		fmt.Println("Query error")
		return -1, err
	}

	row.Scan(&id_member, &pwd)
	fmt.Println(pwd)
	fmt.Println(id_member)
	match, err := CheckHashedPassword(password, pwd)
	if !match {
		fmt.Println("Hash and password does not match.")
		return -1, err
	}
	return id_member, nil
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

func WatchFilm(w http.ResponseWriter, r *http.Request) {
	type FilmHeader struct {
		Judul       string   `json:"judul"`
		TahunRilis  string   `json:"tahunRilis"`
		Sutradara   string   `json:"sutradara"`
		PemainUtama []string `json:"pemainutama"`
		Genre       string   `json:"genre"`
		Sinopsis    string   `json:"sinopsis"`
	}
	var filmHeader FilmHeader
	db := db.ConnectDB()
	vars := mux.Vars(r)

	film_id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}

	//Get from film
	res := db.Debug().Where("id_film = ?", film_id).Find(&models.Film{})
	errRes := res.Error

	//If query error
	w.Header().Set("Content-Type", "application/json")
	if errRes != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}
	//If there are no rows
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "There are no films available",
			Data:    nil,
		})
		return
	}
	//Insert to history
	rows, err := db.Debug().Table("pemains a").Select("c.judul, c.tahun_rilis, c.sutradara, a.nama_pemain, d.jenis_genre, c.sinopsis").Joins("join list_pemains b on a.id_pemain = b.id_pemain").Joins("join films c on b.id_film = c.id_film").Joins("join genres d on c.id_genre = d.id_genre").Where("c.id_film = ? and b.peran = ?", film_id, "Pemain Utama").Rows()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var pemain string
		rows.Scan(&filmHeader.Judul, &filmHeader.TahunRilis, &filmHeader.Sutradara, &pemain, &filmHeader.Genre, &filmHeader.Sinopsis)
		filmHeader.PemainUtama = append(filmHeader.PemainUtama, pemain)
	}

	//Validate from cookies
	status, userID, err := GetIDFromCookies(r)
	if !status && err != nil {
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}

	//Insert to history and check if there is an error
	if err := db.Create(&models.History{
		IdMember:      userID,
		IdFilm:        film_id,
		TanggalNonton: time.Now(),
	}).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}

	//Output JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  http.StatusOK,
		Message: "Enjoy the movie!",
		Data:    filmHeader,
	})

}

func GetIDFromCookies(r *http.Request) (bool, int, error) {
	cookie, err := r.Cookie(tokenName)
	if err != nil {
		return false, -1, err
	}

	accessToken := cookie.Value
	accessClaims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil && !parsedToken.Valid {
		return false, -1, err
	}
	return true, accessClaims.UserID, nil
}
