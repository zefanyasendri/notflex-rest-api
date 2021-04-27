package models

type Film struct {
	IdFilm     int          `form:"idFilm" json:"idFilm" gorm:"primaryKey"`
	Judul      string       `form:"judul" json:"judul"`
	TahunRilis string       `form:"tahunRilis" json:"tahunRilis"`
	Sutradara  string       `form:"sutradara" json:"sutradara"`
	Sinopsis   string       `form:"sinopsis" json:"sinopsis"`
	IdGenre    int          `form:"idGenre" json:"idGenre"`
	ListPemain []ListPemain `gorm:"foreignKey:IdFilm"`
	History    []History    `gorm:"foreignKey:IdFilm"`
}

type FilmResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}
