package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	database "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	pass := r.URL.Query()["password"]
	email := r.URL.Query()["email"]
	var member models.Member
	var response models.MemberResponse

	if err := db.Where("email = ? and password = ?", email[0], pass[0]).First(&member).Error; err != nil {
		log.Print(err)
		response.Status = 400
		response.Message = "Error"
		return
	}

	db.Find(&member)
	generateToken(w, member.Email, member.Password, 0)
	response.Status = 200
	response.Message = "Success Login <WELCOME ADMIN>"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMemberBaseOnEmail(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()

	email := r.URL.Query()["email"]

	var member models.Member
	var members []models.Member
	var response models.MemberResponse

	if err := db.Where("email = ?", email[0]).First(&member).Error; err != nil {
		log.Print(err)
		response.Status = 400
		response.Message = "Error"
		return
	}

	members = append(members, member)

	db.Find(&member)
	response.Status = 200
	response.Message = "Success"
	response.Data = members

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
