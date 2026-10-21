package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/api/middlewares"
	"github.com/sebx771/monitor-bot/internal/config"
	"github.com/sebx771/monitor-bot/internal/status"
)

func SetupRouter(store *status.Store, cfg *config.Config) *gin.Engine {
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := NewHandler(store)

	v1 := r.Group("/api/v1")
	{
		// Endpoints públicos
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		v1.GET("/status", handler.GetStatus)

		// Endpoints protegidos por API Key
		secretKey := cfg.GetGlobalConfig().GetSecretKey()
		protected := v1.Group("")
		protected.Use(middlewares.ApiKeyAuth(secretKey))
		{
			protected.POST("/workers/:id/trigger", handler.TriggerWorker)
		}
	}

	return r
}