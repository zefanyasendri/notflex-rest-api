package models

type History struct {
	IdMember      int    `json:"idMember" gorm:"primaryKey"`
	IdFilm        int    `json:"idFilm"`
	TanggalNonton string `json:"tanggalNonton"`
	Member        Member `gorm:"foreignKey:IdMember"`
}
