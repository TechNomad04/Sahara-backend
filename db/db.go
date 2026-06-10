package db

import (
	"fmt"
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = sqlDB.Ping()
		if err == nil {
			sqlDB.SetMaxOpenConns(25)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)

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