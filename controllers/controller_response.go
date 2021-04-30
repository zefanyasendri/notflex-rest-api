package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
)

// UnAuthorized Response
func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response models.Response
	response.Status = http.StatusUnauthorized
	response.Message = "Unauthorized Access"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Success Response
func sendSuccessResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  http.StatusOK,
		Message: msg,
		Data:    nil,
	})
}

// Error Response
func sendErrorResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.Response{
		Status:  http.StatusInternalServerError,
		Message: msg,
		Data:    nil,
	})
}
