package main

import (
	"log"

	"github.com/Casagrande-Lucas/dnd/config"
	"github.com/Casagrande-Lucas/dnd/infrastructure/db"
	"github.com/Casagrande-Lucas/dnd/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	cfg := config.GetConfig()
	factoryDB := db.GetDBFactory()

	dbConn, err := factoryDB.CreatePostgresConnection("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Error connecting to postgres: %v", err)
	}

	srv := api.NewGinRoutes(r, cfg, dbConn.GetDB())
	if err := srv.StartServer(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
