package models

type Person struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PersonResponse struct {
	Status  int    `form:"status" json:"status"`
	Message string `json:"Login Status "`
}
