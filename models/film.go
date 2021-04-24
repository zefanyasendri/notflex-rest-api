package models

type Film struct {
	IdFilm     int        `json:"idFilm" gorm:"primaryKey"`
	Judul      string     `json:"judul"`
	TahunRilis string     `json:"tahunRilis"`
	Sutradara  string     `json:"sutradara"`
	Sinopsis   string     `json:"sinopsis"`
	IdGenre    int        `json:"idGenre"`
	ListPemain ListPemain `gorm:"foreignKey:IdFilm"`
	History    History    `gorm:"foreignKey:IdFilm"`
}
