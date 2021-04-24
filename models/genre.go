package models

type Genre struct {
	IdGenre    int    `json:"idGenre" gorm:"primaryKey"`
	JenisGenre string `json:"jenisGenre"`
	Film       Film   `gorm:"foreignKey:IdGenre"`
}
