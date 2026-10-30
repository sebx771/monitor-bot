package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/api/handlers"
	"github.com/sebx771/monitor-bot/internal/api/middlewares"
	"github.com/sebx771/monitor-bot/internal/auth"
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
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
			// TODO: add production origin (e.g. https://your-app.vercel.app)
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	globalCfg := cfg.GetGlobalConfig()

	authSvc := auth.NewSessionService(
		globalCfg.GetAdminUserName(),
		globalCfg.GetAdminPasswordHash(),
	)

	statusHandler := handlers.NewStatusHandler(store)
	authHandler := handlers.NewAuthHandler(authSvc)

	v1 := r.Group("/api/v1")
	{
		// ── Public endpoints ──────────────────────────────────────────────
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		v1.GET("/status", statusHandler.GetStatus)
		v1.GET("/auth/me", authHandler.Me)
		v1.POST("/auth/login", authHandler.Login)

		// ── Session-protected endpoints ────────────────────────────────────
		sessionProtected := v1.Group("")
		sessionProtected.Use(middlewares.SessionAuth(authSvc))
		{
			sessionProtected.POST("/auth/logout", authHandler.Logout)
			sessionProtected.POST("/workers/:id/trigger", statusHandler.TriggerWorker)
		}

		// ── Legacy API-key-protected endpoints (kept during migration) ─────
		// Remove after confirming session auth works end-to-end.
		secretKey := globalCfg.GetSecretKey()
		if secretKey != "" {
			apiKeyProtected := v1.Group("")
			apiKeyProtected.Use(middlewares.ApiKeyAuth(secretKey))
			// (no routes currently registered here; preserved for rollback)
			_ = apiKeyProtected
		}
	}

	return r
}