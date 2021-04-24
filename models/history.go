package models

type History struct {
	IdMember      int    `json:"idMember"`
	IdFilm        int    `json:"idFilm"`
	TanggalNonton string `json:"tanggalNonton"`
}
