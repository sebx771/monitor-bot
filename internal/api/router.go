package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/status"
)

func SetupRouter(store *status.Store) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, 
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := NewHandler(store)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		v1.GET("/status", handler.GetStatus)
		v1.POST("/workers/:id/trigger", handler.TriggerWorker)
	}

	return r
}