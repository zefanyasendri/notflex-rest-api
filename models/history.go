package models

type History struct {
	IdMember      int    `json:"idMember"`
	IdFilm        int    `json:"idFilm"`
	TanggalNonton string `json:"tanggalNonton"`
}

type HistoryResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}