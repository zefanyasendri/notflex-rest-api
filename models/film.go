package models

type Film struct {
	IdFilm     int    `json:"idFilm"`
	Judul      string `json:"judul"`
	TahunRilis string `json:"tahunRilis"`
	Sutradara  string `json:"sutradara"`
	Sinopsis   string `json:"sinopsis"`
	IdGenre    int    `json:"idGenre"`
}
