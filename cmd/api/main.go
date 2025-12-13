package main

import (
	"context"
	"log"
	stdhttp "net/http"
	"os"
	"project-manager-dashboard-go/internal/app/usecase/project"
	"time"

	"github.com/joho/godotenv"

	"project-manager-dashboard-go/internal/app"
	"project-manager-dashboard-go/internal/app/usecase/user"
	httpapi "project-manager-dashboard-go/internal/transport/http"
)

func main() {
	_ = godotenv.Load()

	addr := getenv("HTTP_ADDR", ":8081")
	dbURL := getenv("DATABASE_URL", "")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	a, err := app.New(dbURL)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer a.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Ent.Schema.Create(ctx); err != nil {
		log.Fatalf("schema create: %v", err)
	}

	// Users
	userRepo := user.NewUserRepo(a.Ent)
	userUseCase := user.NewUserUsecase(userRepo)
	userHandlers := httpapi.NewUserHandler(userUseCase)

	// Projects
	projectRepo := project.NewEntRepo(a.Ent)
	projectUseCase := project.NewProjectUsecase(projectRepo)
	projectHandlers := httpapi.NewProjectHandler(projectUseCase)

	r := httpapi.NewRouter(userHandlers, projectHandlers)

	log.Printf("HTTP listening on %s", addr)
	log.Fatal(stdhttp.ListenAndServe(addr, r))
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
