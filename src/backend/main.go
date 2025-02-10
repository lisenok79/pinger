package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Request struct {
	ContainerID string            `json:"containerID"`
	IP          map[string]string `json:"ip"`
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Datestamp   time.Time         `json:"datestamp"`
}

type ToFront struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Datestamp time.Time `json:"datestamp"`
}

type DBContainer struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerID string    `gorm:"uniqueIndex;not null"`
	IP          string    `gorm:"type:varchar(255);not null"`
	Status      string    `gorm:"type:varchar(255);not null"`
	Timestamp   time.Time `gorm:"not null"`
	Datestamp   time.Time `gorm:"not null"`
}

type Env struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPort     string
	DBName     string
	DBPassword string
}

func ParseEnv() Env {
	var env Env
	env.Port = os.Getenv("PORT")
	if env.Port == "" {
		log.Fatal("PORT environment variable is required")
		os.Exit(1)
	}
	env.DBHost = os.Getenv("DATABASE_HOST")
	if env.DBHost == "" {
		log.Fatal("DATABASE_HOST environment variable is required")
		os.Exit(1)
	}
	env.DBUser = os.Getenv("DATABASE_USER")
	if env.DBUser == "" {
		log.Fatal("DATABASE_USER environment variable is required")
		os.Exit(1)
	}
	env.DBPort = os.Getenv("DATABASE_PORT")
	if env.DBPort == "" {
		log.Fatal("DATABASE_PORT environment variable is required")
		os.Exit(1)
	}
	env.DBName = os.Getenv("DATABASE_NAME")
	if env.DBName == "" {
		log.Fatal("DATABASE_NAME environment variable is required")
		os.Exit(1)
	}
	env.DBPassword = os.Getenv("DATABASE_PASSWORD")
	if env.DBPassword == "" {
		log.Fatal("DATABASE_PASSWORD environment variable is required")
		os.Exit(1)
	}
	return env
}

func ContainerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Error: wrong http method")
		return
	}
	reqs := []ToFront{}
	conts := []DBContainer{}
	db, err := DbConnect()
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
	reqs := []Request{}
	err = json.Unmarshal(byteReq, &reqs)
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
	err = SaveContainer(db, reqs)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DbConnect() (*gorm.DB, error) {
	env := ParseEnv()
	dsn := fmt.Sprintf("host=%s user=%s port=%s dbname=%s password=%s", env.DBHost, env.DBUser, env.DBPort, env.DBName, env.DBPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&DBContainer{})

	return db, nil
}

func SaveContainer(db *gorm.DB, reqs []Request) error {
	Containers := make([]DBContainer, len(reqs))

	for i, req := range reqs {
		for key, ip := range req.IP {
			Containers[i].IP = key + ", " + ip + "\n"
		}
		Containers[i].ContainerID = req.ContainerID
		Containers[i].Status = req.Status
		Containers[i].Timestamp = req.Timestamp
		Containers[i].Datestamp = req.Datestamp
	}
	for _, Container := range Containers {
		err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "container_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"ip":        gorm.Expr("CASE WHEN ? = 'running' THEN ? ELSE db_containers.ip END", Container.Status, Container.IP),
				"status":    Container.Status,
				"timestamp": Container.Timestamp,
				"datestamp": gorm.Expr("CASE WHEN ? = 'running' THEN ? ELSE db_containers.datestamp END", Container.Status, Container.Datestamp),
			}),
		}).Create(&Container).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	env := ParseEnv()
	mux := http.NewServeMux()
	mux.HandleFunc("/putStatus", PutStatus)
	mux.HandleFunc("/ContainerList", ContainerList)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":"+env.Port, handler)
}
