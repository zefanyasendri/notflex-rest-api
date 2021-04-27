package models

import "time"

type History struct {
	IdMember      int       `json:"idMember" gorm:"primaryKey"`
	IdFilm        int       `json:"idFilm" gorm:"primaryKey"`
	TanggalNonton time.Time `json:"tanggalNonton"`
}
