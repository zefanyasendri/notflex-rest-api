package models

type KartuKredit struct {
	IdMember      int    `json:"idMember" gorm:"primaryKey"`
	NoKartuKredit string `json:"noKartuKredit"`
	MasaBerlaku   string `json:"masaBerlaku"`
	CVC           string `json:"cvc"`
}
