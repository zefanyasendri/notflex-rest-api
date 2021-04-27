package models

import "time"

type History struct {
	IdMember      int       `json:"idMember"`
	IdFilm        int       `json:"idFilm"`
	TanggalNonton time.Time `json:"tanggalNonton"`
}

type HistoryResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}