package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

)

func InitDB() *sql.DB {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Connected to db")
			return db
		}

		log.Println("Waiting for database...")
		log.Println(err)

		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to database")
	return nil
}