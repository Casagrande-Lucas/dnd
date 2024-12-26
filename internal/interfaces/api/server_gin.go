package api

import (
	"github.com/Casagrande-Lucas/dnd/config"
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/controllers"
	persistenceGorm "github.com/Casagrande-Lucas/dnd/internal/domain/race/repositories"
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ginServer struct {
	app    *gin.Engine
	cfg    *config.Config
	dbConn *gorm.DB
}

func NewGinRoutes(app *gin.Engine, cfg *config.Config, dbConn *gorm.DB) Server {
	return &ginServer{
		app:    app,
		cfg:    cfg,
		dbConn: dbConn,
	}
}

func (g *ginServer) ginMode() {
	switch g.cfg.APP.ENV {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "dev":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

func (g *ginServer) getCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     g.cfg.CORS.AllowOrigins,
		AllowMethods:     g.cfg.CORS.AllowMethods,
		AllowHeaders:     g.cfg.CORS.AllowHeaders,
		ExposeHeaders:    g.cfg.CORS.ExposeHeaders,
		AllowCredentials: g.cfg.CORS.AllowCredentials,
	})
}

func (g *ginServer) RegisterServerRoutes() {

	raceRepo := persistenceGorm.NewGormRaceRepository(g.dbConn)
	raceService := services.NewRaceService(raceRepo)
	raceController := controllers.NewRaceControllerGin(raceService)

	g.app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	v1Group := g.app.Group("/api/v1")
	{
		raceV1Group := v1Group.Group("/races")
		{
			raceV1Group.GET("/", raceController.GetAllRaces)
			raceV1Group.GET("/:id", raceController.GetRaceByID)
			raceV1Group.POST("/", raceController.CreateRace)
			raceV1Group.PUT("/:id", raceController.UpdateRace)
			raceV1Group.DELETE("/:id", raceController.DeleteRace)
			raceV1Group.POST("/:id/subraces", raceController.AddSubrace)
			raceV1Group.DELETE("/:id/subraces/:subraceID", raceController.RemoveSubrace)
			raceV1Group.POST("/:id/traits/:traitID", raceController.AddTrait)
			raceV1Group.DELETE("/:id/traits/:traitID", raceController.RemoveTrait)
			raceV1Group.GET("/search", raceController.SearchRaces)
		}
	}
}

func (g *ginServer) StartServer() error {
	g.ginMode()
	g.getCORS()
	g.RegisterServerRoutes()
	if err := g.app.Run(":" + g.cfg.Server.Port); err != nil {
		return err
	}
	return nil
}
