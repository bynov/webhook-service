package main

import (
	"github.com/bynov/webhook-service/internal/config"
	"github.com/bynov/webhook-service/internal/migration"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}

	err = migration.MigratePostgres("file://./migration/postgres", cfg.DatabaseAddr)
	if err != nil {
		panic(err)
	}
}
