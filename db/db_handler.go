package controllers

import (
	"database/sql"
	"log"

	models "github.com/zefanyasendri/TugasKelompok-REST-API-NotFlex/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//sementara pake go-sql dulu, kalau nanti mau pake gorm, tinggal ganti
func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_tubes_notflex?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/db_tubes_notflex?parseTime=True&loc=Asia%2FJakarta"), &gorm.Config{})

	if err != nil {
		panic("Connection Failed")
	}

	db.AutoMigrate(&models.Person{})
	db.AutoMigrate(&models.Member{})
	db.AutoMigrate(&models.KartuKredit{})
	db.AutoMigrate(&models.Genre{})
	db.AutoMigrate(&models.Film{})
	db.AutoMigrate(&models.History{})
	db.AutoMigrate(&models.Pemain{})
	db.AutoMigrate(&models.ListPemain{})

	return db
}
