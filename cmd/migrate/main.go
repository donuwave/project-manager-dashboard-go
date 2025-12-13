package main

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"log"
	"os"

	"project-manager-dashboard-go/ent"
	"project-manager-dashboard-go/ent/migrate"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer drv.Close()

	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	ctx := context.Background()

	if err := client.Schema.Create(
		ctx,
		migrate.WithForeignKeys(true),
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Println("migration complete")
}
