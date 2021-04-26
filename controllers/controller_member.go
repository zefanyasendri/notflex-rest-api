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
