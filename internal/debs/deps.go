package deps

import (
	"log"

	"github.com/bwise1/your_care_api/config"
	"github.com/bwise1/your_care_api/internal/db"
)

type Dependencies struct {
	DB *db.DB
}

func New(cfg *config.Config) *Dependencies {
	database, err := db.New(cfg.Dsn)
	if err != nil {
		log.Panicln("failed to connect to database", "error", err)
	}
	deps := Dependencies{
		DB: database,
	}
	return &deps
}
