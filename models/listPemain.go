package models

type ListPemain struct {
	IdPemain int    `json:"idPemain"`
	IdFilm   int    `json:"idFilm"`
	Peran    string `json:"peran"`
}

type ListPemainResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}