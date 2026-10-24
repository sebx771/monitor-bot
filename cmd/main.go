package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sebx771/monitor-bot/internal/api"
	config "github.com/sebx771/monitor-bot/internal/config"
	"github.com/sebx771/monitor-bot/internal/integrations/aiven"
	"github.com/sebx771/monitor-bot/internal/integrations/aternos"
	"github.com/sebx771/monitor-bot/internal/integrations/supabase"
	"github.com/sebx771/monitor-bot/internal/logger"
	"github.com/sebx771/monitor-bot/internal/status"
	"github.com/sebx771/monitor-bot/internal/worker"
)

const (
	checkInterval = 1440 * time.Minute // Server check frequency
	errCooldown   = 70 * time.Minute   // Wait time if Aternos fails

	aivenInterval = 60 * time.Minute // Aiven services check frequency
	aivenCooldown = 30 * time.Minute // Wait time if the Aiven API fails
)

func main() {
	log := logger.NewLogger("MAIN")
	store := status.NewStore()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("error loading configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	aivenEnabled := cfg.GetAivenConfig().IsEnabled()
	aternosEnabled := cfg.GetAternosConfig().IsEnabled()
	supabaseEnabled := cfg.GetSupaBaseConfig().IsEnabled()

	log.Info("active modules", "|AIVEN|", aivenEnabled, "|ATERNOS|", aternosEnabled, "|SUPABASE|", supabaseEnabled)

	if !aivenEnabled {
		if err := status.StorePausedWorker(store, "aiven", "Aiven"); err != nil {
			log.Error("error storing paused Aiven status", "error", err)
		}
	}
	if !aternosEnabled {
		if err := status.StorePausedWorker(store, "aternos", "Minecraft"); err != nil {
			log.Error("error storing paused Minecraft status", "error", err)
		}
	}
	if !supabaseEnabled {
		if err := status.StorePausedWorker(store, "supabase", "Supabase"); err != nil {
			log.Error("error storing paused Supabase status", "error", err)
		}
	}

	if !aivenEnabled && !aternosEnabled && !supabaseEnabled {
		log.Warn("no module enabled; API will serve paused worker statuses")
	}

	r := api.SetupRouter(store, cfg)
	srv := &http.Server{Addr: ":8080", Handler: r}
	// goroutines for the server run
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error in HTTP server", "error", err)
		}
	}()

	var wg sync.WaitGroup

	// goroutine for the server stop (if context is done, stop the server)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("error shutting down server gracefully", "error", err)
		}
	}()

	// Aternos (Minecraft) bot — only if enabled
	if aternosEnabled {
		at := cfg.GetAternosConfig()

		aternosTask := aternos.NewTask(at, store)
		w, err := worker.New(checkInterval, errCooldown, aternosTask)
		if err != nil {
			log.Error("error initializing Aternos worker", "error", err)
			os.Exit(1)
		}

		startWorker(ctx, &wg, log, "Minecraft", "aternos", w, store)
	}

	// Aiven monitor — only if enabled
	if aivenEnabled {
		aivenTask := aiven.NewTask(cfg.GetAivenConfig().GetCredentials(), store)

		wAiven, err := worker.New(aivenInterval, aivenCooldown, aivenTask)
		if err != nil {
			log.Error("error initializing Aiven worker", "error", err)
			os.Exit(1)
		}

		startWorker(ctx, &wg, log, "Aiven", "aiven", wAiven, store)
	}

	if supabaseEnabled {
		supabaseTask := supabase.NewTask(cfg.GetSupaBaseConfig().GetToken(), store)

		wSupaBase, err := worker.New(aivenInterval, aivenCooldown, supabaseTask)
		if err != nil {
			log.Error("error initializing Supabase worker", "error", err)
			os.Exit(1)
		}

		startWorker(ctx, &wg, log, "Supabase", "supabase", wSupaBase, store)
	}

	wg.Wait()

	log.Info("application finished successfully")
}

func startWorker(ctx context.Context, wg *sync.WaitGroup, log *logger.Logger, name, id string, w *worker.Worker, store *status.Store) {
	store.RegisterWorker(id, w)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := w.Run(ctx); err != nil {
			log.Error(name+" worker finished with error", "error", err)
		}
	}()
}
