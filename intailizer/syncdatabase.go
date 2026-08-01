package intailizer

import (
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Syncdatabase() {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("error setting migration dialect: %v", err)
	}
	if err := goose.Up(DB, "migrations"); err != nil {
		log.Fatalf("error applying database migrations: %v", err)
	}
	log.Println("database migrations are up to date")
}
