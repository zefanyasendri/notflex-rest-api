package models

type KartuKredit struct {
	NoKartuKredit string `json:"noKartuKredit"  gorm:"primaryKey"`
	MasaBerlaku   string `json:"masaBerlaku"`
	CVC           string `json:"cvc"`
}
