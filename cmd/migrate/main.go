package main

import (
	"fmt"
	"log"
	"os"

	"github.com/VarunSuddala/olx-api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()
	m, err := migrate.New("file://migrations", cfg.DBURL)
	if err != nil {
		log.Fatalf("migrate.new %v", err)
	}
	if len(os.Args) < 2 {
		log.Fatalf("usage : make migrate <up | down>")
	}
	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("migrate.up %v", err)
		}

		log.Println("migrating up")
	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("migrate.down %v", err)
		}
		log.Println("migrating down")
	default:
		log.Fatalln("unkown command")
	}
	fmt.Println("running migration")
}
