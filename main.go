package main

import (
	"github.com/bynov/webhook-service/internal/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}
}
