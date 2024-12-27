// @title D&D 5e API
// @version 1.0
// @description API D&D 5e card game.

// @contact.name Lucas Casagrande
// @contact.url http://github.com/Casagrande-Lucas
// @contact.email dev.casagrande@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"log"

	"github.com/Casagrande-Lucas/dnd/config"
	_ "github.com/Casagrande-Lucas/dnd/docs"
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
