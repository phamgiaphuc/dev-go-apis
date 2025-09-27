package database

import (
	"dev-go-apis/internal/lib"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func NewDatabaseClient() *sqlx.DB {
	db, err := sqlx.Connect("postgres", lib.DATABASE_URL)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("🎉 Database is connected")

	return db
}
