package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	DB_URL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Dockerized Go App!")
		fmt.Println("Request processed")
	})

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			http.Error(w, "DB Connection Failed", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Connected to DB!")
	})

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
