package main

import (
	"log"
	"net/http"

	"github.com/fekom/udt/config"
	"github.com/fekom/udt/pkg/database"
	"github.com/joho/godotenv"
)

type App struct {
	addr string
	db   *database.DB
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, reading from environment")
	}

	cfg := config.Load()

	log.Printf("connecting to DB: uri=%s database=%s\n", cfg.DBconfig.URI, cfg.DBconfig.Database)

	db, err := database.New(cfg.DBconfig.URI, cfg.DBconfig.Database)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	app := App{
		addr: cfg.Addr,
		db:   db,
	}
	app.start(cfg.Addr)
}

func (app *App) start(addr string) {
	mux := http.NewServeMux()
	SetupRoutes(mux, app.db)

	log.Printf("server running on %s\n", app.addr)
	if err := http.ListenAndServe(app.addr, mux); err != nil {
		log.Fatal("server error", err)
	}
}
