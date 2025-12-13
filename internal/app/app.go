package app

import (
	"database/sql"

	"project-manager-dashboard-go/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Ent *ent.Client
	DB  *sql.DB
}

func New(databaseURL string) (*App, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	return &App{Ent: client, DB: db}, nil
}

func (a *App) Close() {
	if a.Ent != nil {
		_ = a.Ent.Close()
	}
	if a.DB != nil {
		_ = a.DB.Close()
	}
}
