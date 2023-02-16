package db

import (
	"fmt"

	config "github.com/Asliddin3/cykel-omni/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDb(cfg config.Config) (*sqlx.DB, error) {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)
	connDb, err := sqlx.Connect("postgres", psqlString)
	if err != nil {
		return nil, err
	}
	return connDb, nil
}

func ConnectToDbForSuite(cfg config.Config) (*sqlx.DB, func()) {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)
	connDb, err := sqlx.Connect("postgres", psqlString)
	if err != nil {
		return nil, func() {}
	}
	cleandUpfunc := func() {
		connDb.Close()
	}
	return connDb, cleandUpfunc
}
