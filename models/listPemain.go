package models

type ListPemain struct {
	IdPemain int    `json:"idPemain" gorm:"primaryKey"`
	IdFilm   int    `json:"idFilm" gorm:"primaryKey"`
	Peran    string `json:"peran"`
}

type ListPemainResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}