package models

type ListPemain struct {
	IdPemain int    `json:"idPemain" gorm:"primaryKey"`
	IdFilm   int    `json:"idFilm"`
	Peran    string `json:"peran"`
}
