package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

	res1, err := CheckLogin(email, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Email/password seems to be incorrect. Please try again.",
			Data:    nil,
		})
		return
	}

	res2, err := CheckSuspended(email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErrorResponse(w, "Something went wrong. Please try again")
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
	row := db.Table("members").Where("email = ?", email).Select("id_member", "password").Row()
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
	//Validate from cookies
	status, userID, err := GetIDFromCookies(r)
	if !status && err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
		return
	}

	if !CheckSubscribe(userID) {
		sendErrorResponse(w, "Looks like you're not subscribed yet.")
		return
	}

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
		sendErrorResponse(w, "Something went wrong. Please try again")
		return
	}

	//Get from film
	res := db.Where("id_film = ?", film_id).Find(&models.Film{})
	errRes := res.Error

	//If query error
	w.Header().Set("Content-Type", "application/json")
	if errRes != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
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
	rows, err := db.Table("pemains a").Select("c.judul, c.tahun_rilis, c.sutradara, a.nama_pemain, d.jenis_genre, c.sinopsis").Joins("join list_pemains b on a.id_pemain = b.id_pemain").Joins("join films c on b.id_film = c.id_film").Joins("join genres d on c.id_genre = d.id_genre").Where("c.id_film = ? and b.peran = ?", film_id, "Pemain Utama").Rows()
	if err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var pemain string
		rows.Scan(&filmHeader.Judul, &filmHeader.TahunRilis, &filmHeader.Sutradara, &pemain, &filmHeader.Genre, &filmHeader.Sinopsis)
		filmHeader.PemainUtama = append(filmHeader.PemainUtama, pemain)
	}

	//Insert to history and check if there is an error
	if err := db.Create(&models.History{
		IdMember:      userID,
		IdFilm:        film_id,
		TanggalNonton: time.Now().Format("02/01/2006"),
	}).Error; err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
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

func Subscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get card value
	kartu_kredit := r.FormValue("kartu_kredit")
	cvc := r.FormValue("cvc")
	masa_berlaku := r.FormValue("masa_berlaku")
	paket_pilihan := r.FormValue("paket_pilihan")

	//ambil id user dari cookies
	status, userID, err := GetIDFromCookies(r)
	if !status && err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again. ")
		return
	}

	//Check status kalo udah subscribe ngapain lagi subscribe ulang?
	if CheckSubscribe(userID) {
		sendErrorResponse(w, "You're already subscribed")
		return
	}

	db := db.ConnectDB()

	//Mencari kartu kredit yang mungkin pernah dimasukkan oleh user
	var db_res models.KartuKredit
	query := db.Where("id_member = ?", userID).Find(&db_res)

	if err := query.Error; err != nil {
		sendErrorResponse(w, "Something went wrong please try again")
		fmt.Println(err.Error())
		return
	}

	//Jika belum ada
	if query.RowsAffected == 0 {
		//Insert data kartu baru
		if err := db.Create(&models.KartuKredit{
			IdMember:      userID,
			NoKartuKredit: kartu_kredit,
			MasaBerlaku:   masa_berlaku,
			CVC:           cvc,
		}).Error; err != nil {
			sendErrorResponse(w, "Something went wrong. Please try again")
			return
		}
	} else {
		//Check kalau data yang dimasukkan sesuai dengan di database atau tidak
		if db_res.NoKartuKredit != kartu_kredit || db_res.MasaBerlaku != masa_berlaku || db_res.CVC != cvc {
			sendErrorResponse(w, "One of your credentials are incorrect.")
			return
		}
	}

	// Kalo kartu kredit expired tidak bisa membayar
	t, err := time.Parse("02/01/2006", masa_berlaku)
	if err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
		return
	}
	if time.Now().After(t) {
		sendErrorResponse(w, "Your card has expired")
		return
	}

	//Cek pengisian paket apakah benar atau salah
	switch paket_pilihan {
	case "basic":
		paket_pilihan = "Subscribed - Basic"
	case "Basic":
		paket_pilihan = "Subscribed - Basic"
	case "premium":
		paket_pilihan = "Subscribed - Premium"
	case "Premium":
		paket_pilihan = "Subscribed - Premium"
	default:
		sendErrorResponse(w, "Subscription not found")
		return
	}

	// jika semua sesuai maka akan diupdate
	res1 := db.Model(&models.Member{}).Where("id_member = ?", userID).Updates(map[string]interface{}{
		"status_akun":        paket_pilihan,
		"subscription_until": time.Now().AddDate(0, 0, 30).Format("02/01/2006"),
	})

	if err := res1.Error; err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
		return
	}
	sendSuccessResponse(w, "Your subscription has been purchased until "+time.Now().AddDate(0, 0, 30).Format("02 January, 2006")+". Enjoy your Notflex!")
}

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get id from cookies
	status, userID, err := GetIDFromCookies(r)
	if !status && err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again. ")
		return
	}

	//Kebalikan subscribe kalo ini cek apakah unsubscribe atau tidak.
	if !CheckSubscribe(userID) {
		sendErrorResponse(w, "You're unsubscribed")
		return
	}

	//jika statusnya subscribe maka tabel diupdate
	db := db.ConnectDB()
	res := db.Model(&models.Member{}).Where("id_member = ?", userID).Updates(map[string]interface{}{
		"status_akun":        "On Hold",
		"subscription_until": time.Now().Format("02/01/2006"),
	})
	if err := res.Error; err != nil {
		sendErrorResponse(w, "Something went wrong. Please try again")
		return
	}
	sendSuccessResponse(w, "Unsubscribed! Sad to see you go.")
}

func CheckSubscribe(id_member int) bool {
	var status string
	db := db.ConnectDB()

	row := db.Table("members").Where("id_member = ?", id_member).Select("status_akun").Find(&status)
	if err := row.Error; err != nil {
		return false
	}

	if status == "Subscribed - Basic" || status == "Subscribed - Premium" {
		return true
	}
	return false
}

func GetWatchHistory(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Judul string    `json:"judul"`
		Waktu time.Time `json:"waktu"`
	}
	//Validate from cookies
	status, id_member, err := GetIDFromCookies(r)
	if !status && err != nil {
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}
	db := db.ConnectDB()
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

func Register(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	var newmember models.Member
	json.Unmarshal(body, &newmember)
	var emailscan string
	query, _ := db.Debug().Table("members").Select("email").Where("email = ?", newmember.Email).Rows()
	fmt.Println(newmember.Email)

	for query.Next() {
		query.Scan(&emailscan)
	}
	fmt.Println(emailscan)
	if len(emailscan) == 0 {
		db.Save(&newmember)
		response := models.FilmResponse{Status: 200, Data: newmember, Message: "WELCOME ABOARD!!!"}
		result, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	} else {
		response := models.FilmResponse{Status: 400, Message: "Email telah diambil alih"}
		result, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
