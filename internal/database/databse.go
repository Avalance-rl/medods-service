package database

import (
	"log"
	"medods-service/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)
var (
	DB *sqlx.DB
)
func Init() {
	DB = mustCreate(config.CFG)
}

func mustCreate(cfg config.Config) *sqlx.DB {
	dbConnArg := cfg.Database.Address
	db, err := sqlx.Connect("pgx", dbConnArg)
	if err != nil {
		log.Fatalf("mustCreate() Error: %v", err)
	}

	return db
}

