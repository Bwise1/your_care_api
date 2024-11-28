package deps

import (
	"log"

	"github.com/bwise1/your_care_api/config"
	"github.com/bwise1/your_care_api/internal/db"
	smtp "github.com/bwise1/your_care_api/util/email"
)

type Dependencies struct {
	DB     *db.DB
	Mailer *smtp.Mailer
}

func New(cfg *config.Config) *Dependencies {
	database, err := db.New(cfg.Dsn)
	mailer := smtp.NewMailer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPassword, cfg.SmtpFrom)

	if err != nil {
		log.Panicln("failed to connect to database", "error", err)
	}
	deps := Dependencies{
		DB:     database,
		Mailer: mailer,
	}
	return &deps
}
