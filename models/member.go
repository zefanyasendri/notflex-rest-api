package models

type Member struct {
	Person        `json:"person"`
	IdMember      int         `json:"idMember" gorm:"primaryKey"`
	NamaLengkap   string      `json:"namaLengkap"`
	TanggalLahir  string      `json:"tanggalLahir"`
	JenisKelamin  string      `json:"jenisKelamin"`
	AsalNegara    string      `json:"asalNegara"`
	StatusAkun    string      `json:"statusAkun"`
	NoKartuKredit string      `json:"noKartuKredit"`
	KartuKredit   KartuKredit `gorm:"foreignKey:NoKartuKredit"`
}
