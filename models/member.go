package models

type Member struct {

	Person            `json:"person"`
	IdMember          int         `json:"idMember" gorm:"primaryKey"`
	NamaLengkap       string      `json:"namaLengkap"`
	TanggalLahir      string      `json:"tanggalLahir"`
	JenisKelamin      string      `json:"jenisKelamin"`
	AsalNegara        string      `json:"asalNegara"`
	StatusAkun        string      `json:"statusAkun"`
	SubscriptionUntil string      `json:"subscriptionUntil"`
	History           History     `gorm:"foreignKey:IdMember"`
	KartuKredit       KartuKredit `gorm:"foreignKey:IdMember"`
}
