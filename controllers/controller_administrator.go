package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	database "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()

	pass := r.URL.Query()["password"]
	email := r.URL.Query()["email"]

	query := "SELECT * FROM person"
	if email != nil && pass != nil {
		query += " WHERE password ='" + pass[0] + "' and email = '" + email[0] + "'"
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Print(err)
	}

	var person models.Person
	var persons []models.Person
	for rows.Next() {
		if err := rows.Scan(&person.Email, &person.Password); err != nil {
			log.Print(err.Error())
		} else {
			persons = append(persons, person)
		}
	}

	var response models.PersonResponse
	if err == nil {
		generateToken(w, person.Email, person.Password, 0)
		response.Status = 200
		response.Message = "Success Login <WELCOME>"
	} else {
		response.Status = 400
		response.Message = "Error"
	}

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
