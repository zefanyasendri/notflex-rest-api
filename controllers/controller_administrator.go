package controllers

import (
	"encoding/json"
	"log"
	"net/http"

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
	generateToken(w, person.Email, person.Password, 0)
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
