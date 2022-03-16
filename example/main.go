package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	sql_console "github.com/chandlerswift/go-sql-console"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if _, err := os.Stat("./example.db"); errors.Is(err, os.ErrNotExist) {
		log.Println("Example database doesn't exist. Retrieving from GitHub.")
		downloadURL := "https://github.com/lerocha/chinook-database/raw/master/ChinookDatabase/DataSources/Chinook_Sqlite.sqlite"
		response, err := http.Get(downloadURL)
		if err != nil {
			log.Fatalf("Error downloading chinook database: %v", err.Error())
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			log.Fatalf("Unexpected response code %v downloading chinook database", response.StatusCode)
		}
		file, err := os.Create("./example.db")
		if err != nil {
			log.Fatalf("Error creating ./example.db: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("Error writing to ./example.db: %v", err)
		}
	}
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatalf("Error opening ./example.db: %v", err)
	}
	defer db.Close()

	http.Handle("/", sql_console.Handler{
		DB: db,
	})
	log.Printf("Serving...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
