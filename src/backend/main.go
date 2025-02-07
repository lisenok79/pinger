package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBContainer struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `gorm:"uniqueIndex;not null" json:"ip"`
	Status    string    `gorm:"type:varchar(255);not null" json:"status"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	Datestamp string    `gorm:"type:varchar(255)" json:"datestamp"`
}

func PutStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Wrong method!")
		return
	}

	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	dbCont := DBContainer{}
	err = json.Unmarshal(byteReq, &dbCont)
	if err != nil {
		log.Println(err)
		return
	}

	db, err := DbConnect()
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
		return
	}
	defer sqlDB.Close()
	err = SaveContainer(db, dbCont)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DbConnect() (*gorm.DB, error) {
	dsn := "host=postgres user=myuser port=5432 dbname=mydatabase password=''"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&DBContainer{})
	
	return db, nil
}

func SaveContainer(db *gorm.DB, newContainer DBContainer) error {
	// container := DBContainer{
	// 	IP:        "0.0.0.0",
	// 	Status:    "DOWN",
	// 	Timestamp: time.Now(),
	// 	Datestamp: time.Now().String(),
	// }
	// sl, _ := json.Marshal(container)
	// log.Println(string(sl))
	// time.Sleep(2 * time.Second)
	// newContainer := DBContainer{
	// 	IP:        "0.0.0.0",
	// 	Status:    "OK",
	// 	Timestamp: time.Now(),
	// 	Datestamp: time.Now().String(),
	// }
	// res := db.Create(&container)
	// if res.Error != nil {
	// 	return res.Error
	// }
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "ip"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":    newContainer.Status,
			"timestamp": newContainer.Timestamp,
			"datestamp": gorm.Expr("CASE WHEN ? = 'OK' THEN ? ELSE db_containers.datestamp END", newContainer.Status, newContainer.Datestamp),
		}),
	}).Create(&newContainer).Error
	if err != nil {
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/putStatus", PutStatus)
	DbConnect()
	http.ListenAndServe(":8080", nil)
}
