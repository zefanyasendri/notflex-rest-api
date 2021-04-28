package models

type Member struct {
	Email         string      `json:"email"`
	Password      string      `json:"password"`
	IdMember      int         `json:"idMember" gorm:"primaryKey"`
	NamaLengkap   string      `json:"namaLengkap"`
	TanggalLahir  string      `json:"tanggalLahir"`
	JenisKelamin  string      `json:"jenisKelamin"`
	AsalNegara    string      `json:"asalNegara"`
	StatusAkun    string      `json:"statusAkun"`
	NoKartuKredit string      `json:"noKartuKredit" gorm:"type:varchar(191)"`
	History       History     `gorm:"foreignKey:IdMember"`
	KartuKredit   KartuKredit `gorm:"foreignKey:NoKartuKredit"`
}

type MemberResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}
