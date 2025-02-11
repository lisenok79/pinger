package handlers

import (
	"backend/database"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type ToFront struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Datestamp time.Time `json:"datestamp"`
}

func ContainerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Error: wrong http method")
		return
	}
	reqs := []ToFront{}
	conts := []database.DBContainer{}
	db, err := database.DbConnect()
	if err != nil {
		log.Println(err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
		return
	}
	defer sqlDB.Close()
	query := db.Order("timestamp ASC").Find(&conts)
	if query.Error != nil {
		log.Println(query.Error)
		return
	}

	for i := range conts {
		reqs = append(reqs, ToFront{
			IP:        conts[i].IP,
			Status:    conts[i].Status,
			Timestamp: conts[i].Timestamp,
			Datestamp: conts[i].Datestamp,
		})
	}

	json, err := json.Marshal(reqs)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(json)
	w.WriteHeader(http.StatusOK)
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
	reqs := []database.Request{}
	err = json.Unmarshal(byteReq, &reqs)
	if err != nil {
		log.Println(err)
		return
	}

	db, err := database.DbConnect()
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
	err = database.SaveContainer(db, reqs)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
