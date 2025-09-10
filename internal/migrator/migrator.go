package migrator

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

func Migrate(dbURL string) error {
	db, err := sql.Open("pgx", dbURL)
	defer func() {
		err = db.Close()
		if err != nil {
			log.Errorf("could not close db connection:%s", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("could not run migration")
	}
	if err = goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("could not set dialect. err:%v", err)
	}
	if err = goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("could not run migration. err: %v", err)
	}
	return nil
}
