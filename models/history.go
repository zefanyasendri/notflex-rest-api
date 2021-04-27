package models

import "time"

type History struct {
	IdMember      int       `json:"idMember"`
	IdFilm        int       `json:"idFilm"`
	TanggalNonton time.Time `json:"tanggalNonton"`
}
