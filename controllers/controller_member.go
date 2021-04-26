package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	database "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/db"
	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

func LoginMember(w http.ResponseWriter, r *http.Request) {
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
	generateToken(w, member.Email, member.Password, 1)
	response.Status = 200
	response.Message = "Success Login <WELCOME>"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	body, _ := ioutil.ReadAll(r.Body)
	var member models.Member
	json.Unmarshal(body, &member)
	db.Create(&member)

	response := models.MemberResponse{Status: 200, Message: "WELCOME ABOARD!!"}
	result, err := json.Marshal(response)

	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func SignOut(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	var response models.MemberResponse
	response.Status = 200
	response.Message = "Logout Success"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
