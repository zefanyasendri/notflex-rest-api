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

// 1 - User Registration
func Register(w http.ResponseWriter, r *http.Request) {
	type hasil struct {
		NamaLengkap  string `json:"namaLengkap"`
		TanggalLahir string `json:"tanggalLahir"`
		JenisKelamin string `json:"jenisKelamin"`
		AsalNegara   string `json:"asalNegara"`
		Email        string `json:"email"`
		Password     string `json:"password"`
	}

	db := db.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	var newmember models.Member
	var data hasil
	var emailscan string
	json.Unmarshal(body, &data)
	query, _ := db.Debug().Table("members").Select("email").Where("email = ?", data.Email).Rows()
	fmt.Println(data.Email)
	for query.Next() {
		query.Scan(&emailscan)
	}
	newmember.NamaLengkap = data.NamaLengkap
	newmember.TanggalLahir = data.TanggalLahir
	newmember.JenisKelamin = data.JenisKelamin
	newmember.AsalNegara = data.AsalNegara
	newmember.Email = data.Email
	newmember.StatusAkun = "AKTIF"
	fmt.Println(emailscan)
	if len(emailscan) == 0 {
		newmember.Password, _ = HashPassword(data.Password)
		db.Save(&newmember)
		response := models.FilmResponse{Status: 200, Data: data, Message: "WELCOME ABOARD!!!"}
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

// 2 - User Sign-in dan Sign-Out
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

// 2 - User Sign-in dan Sign-Out
func SignOut(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	var response models.Response
	response.Status = 200
	response.Message = "SignOut Success"

	w.Header().Set("Content-Type", "application/json")
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

// 3 - Update Profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	//Validate from cookies
	status, id_member, err := GetIDFromCookies(r)
	if !status && err != nil {
		sendErrorResponse(w, "Something went Wrong, Please try again!")
		return
	}

	db := db.ConnectDB()

	body, _ := ioutil.ReadAll(r.Body)

	type Response struct {
		NamaLengkap       string `json:"namaLengkap"`
		TanggalLahir      string `json:"tanggalLahir"`
		JenisKelamin      string `json:"jenisKelamin"`
		AsalNegara        string `json:"asalNegara"`
		StatusAkun        string `json:"statusAkun"`
		SubscriptionUntil string `json:"subscriptionUntil"`
	}

	var res Response
	var profileUpdates models.Member
	json.Unmarshal(body, &profileUpdates)

	var member models.Member
	db.Model(&member).Where("id_member = ?", id_member).Updates(profileUpdates)
	db.Table("members").Select("id_member, nama_lengkap, tanggal_lahir, jenis_kelamin, asal_negara, status_akun, subscription_until").Where("id_member = ?", id_member).Find(&res)

	response := models.Response{Status: 200, Data: res, Message: "Member Data Updated"}
	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// 4 - Mencari data film sesuai keywords yang diinputkan
func GetFilmByKeywords(w http.ResponseWriter, r *http.Request) {
	//Validate from cookies
	status, _, err := GetIDFromCookies(r)
	if !status && err != nil {
		json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong please try again",
			Data:    nil,
		})
		return
	}

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
	keywordfilm := vars["keywords"]

	var hasil result
	var hasils []result

	query_film, _ := db.Debug().Table("films").Select("films.id_film, films.judul, films.tahun_rilis, films.sutradara, films.sinopsis, films.id_genre").Joins("LEFT JOIN genres ON films.id_genre = genres.id_genre LEFT JOIN list_pemains ON films.id_film = list_pemains.id_film LEFT JOIN pemains ON list_pemains.id_pemain = pemains.id_pemain").Where("films.judul LIKE ? OR films.sutradara LIKE ? OR films.tahun_rilis LIKE ? OR films.sinopsis LIKE ? OR genres.jenis_genre LIKE ? OR pemains.nama_pemain LIKE ?", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%", "%"+keywordfilm+"%").Rows()

	defer query_film.Close()

	for query_film.Next() {

		db.ScanRows(query_film, &hasil)

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

// 5 - Berlangganan
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
	t, err := time.Parse("01/2006", masa_berlaku)
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

// 6 - Berhenti Berlangganan
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

// 7 - "Menonton" Film
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
	rows, err := db.Table("pemains a").Select("c.judul, c.tahun_rilis, c.sutradara, a.nama_pemain, d.jenis_genre, c.sinopsis").Joins("join list_pemains b on a.id_pemain = b.id_pemain").Joins("join films c on b.id_film = c.id_film").Joins("join genres d on c.id_genre = d.id_genre").Where("c.id_film = ?", film_id).Rows()
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

// 8 - Melihat Riwayat Film
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
