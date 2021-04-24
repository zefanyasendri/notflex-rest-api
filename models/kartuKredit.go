package models

type KartuKredit struct {
	NoKartuKredit string `json:"noKartuKredit"`
	MasaBerlaku   string `json:"masaBerlaku"`
	CVC           string `json:"cvc"`
}
