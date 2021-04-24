package models

type Pemain struct {
	IdPemain   int        `json:"idPemain" gorm:"primaryKey"`
	NamaPemain string     `json:"namaPemain"`
	ListPemain ListPemain `gorm:"foreignKey:IdPemain"`
}
