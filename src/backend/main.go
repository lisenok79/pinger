package main

import (
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBContainer struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `gorm:"uniqueIndex;not null"`
	Status    string    `gorm:"type:varchar(255);not null"`
	Timestamp time.Time `gorm:"not null"`
	Datestamp string    `gorm:"type:varchar(255)"`
}

func DbConnect() *gorm.DB {

	dsn := "host=localhost user=myuser port=5433 dbname=mydatabase sslmode=disable password=''"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&DBContainer{})
	SaveContainer(db)
	return db
}

func SaveContainer(db *gorm.DB) error {
	container := DBContainer{
		IP:        "0.0.0.0",
		Status:    "DOWN",
		Timestamp: time.Now(),
		Datestamp: time.Now().String(),
	}
	time.Sleep(2 * time.Second)
	newContainer := DBContainer{
		IP:        "0.0.0.0",
		Status:    "OK",
		Timestamp: time.Now(),
		Datestamp: time.Now().String(),
	}
	res := db.Create(&container)
	if res.Error != nil {
		return res.Error
	}
	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "ip"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":    newContainer.Status,
			"timestamp": newContainer.Timestamp,
			"datestamp": gorm.Expr("CASE WHEN ? = 'OK' THEN ? END", newContainer.Status, newContainer.Datestamp),
		}),
	}).Create(&newContainer)
	return nil
}

func main() {
	// http.HandleFunc("/", pingFunc)
	DbConnect()
	http.ListenAndServe(":8080", nil)
}
