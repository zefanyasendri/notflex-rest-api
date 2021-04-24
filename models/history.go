package models

type History struct {
	IdMember      int    `json:"idMember" gorm:"primaryKey"`
	IdFilm        int    `json:"idFilm" gorm:"primaryKey"`
	TanggalNonton string `json:"tanggalNonton"`
}
