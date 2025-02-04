package main

import (
	"log"
	"net/http"
)

func pingFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("HUI OGROMNIY")
}

func main() {
	http.HandleFunc("/", pingFunc)
	http.ListenAndServe(":8080", nil)
}