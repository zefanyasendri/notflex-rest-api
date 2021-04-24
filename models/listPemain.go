package models

type ListPemain struct {
	IdPemain int    `json:"idPemain" gorm:"primaryKey"`
	IdFilm   int    `json:"idFilm" gorm:"primaryKey"`
	Peran    string `json:"peran"`
}
