package models

type Member struct {
	Person        `json:"person"`
	IdMember      int    `json:"idMember"`
	NamaLengkap   string `json:"namaLengkap"`
	TanggalLahir  string `json:"tanggalLahir"`
	JenisKelamin  string `json:"jenisKelamin"`
	AsalNegara    string `json:"asalNegara"`
	StatusAkun    string `json:"statusAkun"`
	NoKartuKredit string `json:"noKartuKredit"`
}

type MemberResponse struct {
	Status  int      `json:"Status"`
	Message string   `json:"Message"`
	Data    []Member `json:"Data"`
}
