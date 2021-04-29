package models

type Pemain struct {
	IdPemain   int          `json:"idPemain" gorm:"primaryKey"`
	NamaPemain string       `json:"namaPemain"`
	ListPemain []ListPemain `gorm:"foreignKey:IdPemain"`
}

type PemainResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}