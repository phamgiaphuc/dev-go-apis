package database

import (
	"context"
	"database/sql"
	"dev-go-apis/internal/database/migration"
	"dev-go-apis/internal/lib"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"

	_ "github.com/lib/pq"
)

func MigrateDB() {
	log.Printf("🎉 Migrating database")

	ctx := context.Background()
	db, err := sql.Open("postgres", lib.DATABASE_URL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	provider, err := goose.NewProvider(database.DialectPostgres, db, migration.Embed)
	if err != nil {
		log.Fatal(err)
	}

	results, err := provider.Up(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		log.Printf("🎉 %-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}
}
