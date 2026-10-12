package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sebx771/monitor-bot/internal/api"
	config "github.com/sebx771/monitor-bot/internal/config"
	"github.com/sebx771/monitor-bot/internal/integrations/aiven"
	"github.com/sebx771/monitor-bot/internal/logger"
	"github.com/sebx771/monitor-bot/internal/status"

	"github.com/sebx771/monitor-bot/internal/integrations/aternos/adapters"
	"github.com/sebx771/monitor-bot/internal/integrations/aternos/minecraft"
	service "github.com/sebx771/monitor-bot/internal/integrations/aternos/services"
	"github.com/sebx771/monitor-bot/internal/integrations/supabase"

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

	if !aivenEnabled && !aternosEnabled && !supabaseEnabled {
		log.Warn("no module enabled, no workers to run")
		return
	}

	r := api.SetupRouter(store)
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Error("error al iniciar servidor", "error", err)
		}
	}()

	var wg sync.WaitGroup

	// Aternos (Minecraft) bot — only if enabled
	if aternosEnabled {
		at := cfg.GetAternosConfig()

		checker := minecraft.NewChecker(at.GetHost(), at.GetPort())
		gistClient := adapters.NewGitHubGistClient(at.GetGithubToken(), at.GetGistID())
		botService := service.NewBotService(checker, gistClient,
			at.GetStoragePath(), at.GetServerID(), at.IsHeadless(), at.GetHost())

		aternosTask := buildAternosTask(*botService, store)
		w, err := worker.New(checkInterval, errCooldown, aternosTask)
		if err != nil {
			log.Error("error initializing Aternos worker", "error", err)
			os.Exit(1)
		}

		startWorker(ctx, &wg, log, "Minecraft", w)
	}

	// Aiven monitor — only if enabled
	if aivenEnabled {
		aivenTask := buildAivenTask(cfg.GetAivenConfig().GetCredentials(), store)

		wAiven, err := worker.New(aivenInterval, aivenCooldown, aivenTask)
		if err != nil {
			log.Error("error initializing Aiven worker", "error", err)
			os.Exit(1)
		}

		startWorker(ctx, &wg, log, "Aiven", wAiven)
	}

	if supabaseEnabled {
		supabaseTask := buildSupaBaseTask(cfg.GetSupaBaseConfig().GetToken(), store)

		wSupaBase, err := worker.New(aivenInterval, aivenCooldown, supabaseTask)
		if err != nil {
			log.Error("error initializing Supabase worker", "error", err)
			os.Exit(1)
		}

		startWorker(ctx, &wg, log, "Supabase", wSupaBase)
	}

	wg.Wait()

	log.Info("application finished successfully")
}

func startWorker(ctx context.Context, wg *sync.WaitGroup, log *logger.Logger, name string, w *worker.Worker) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := w.Run(ctx); err != nil {
			log.Error(name+" worker finished with error", "error", err)
		}
	}()
}

// buildAivenTask runs a checker per credential. An API failure does not abort
// the others; if at least one fails, it returns a combined error to trigger
// the worker's cooldown.
func buildAivenTask(credentials []config.Credential, store *status.Store) worker.Task {
	return func(ctx context.Context) error {
		var allServices []status.MonitoredService
		var errs []error
		successfulProjects := 0

		for _, cred := range credentials {
			checker := aiven.NewChecker(aiven.NewClient(cred.Token), cred.Project)

			services, err := checker.Check()
			if err != nil {
				errs = append(errs, fmt.Errorf("project %s: %w", cred.Project, err))
				continue
			}
			successfulProjects++

			for _, s := range services {
				allServices = append(allServices, status.MonitoredService{
					Target: s.Name,
					Scope:  cred.Project,
					State:  s.State,
				})
			}
		}

		var ws status.WorkerStatus

		switch {
		case len(errs) == 0:
			ws, _ = status.NewWorkerStatus("aiven", "Aiven", status.Operational, allServices, "")
		case successfulProjects == 0:
			ws, _ = status.NewWorkerStatus("aiven", "Aiven", status.Error, allServices, "All projects failed")
		default:
			ws, _ = status.NewWorkerStatus("aiven", "Aiven", status.Degraded, allServices, "Some projects failed to check")
		}

		if store != nil {
			store.Set(ws)
		}

		// Partial failures mark the worker as degraded without triggering the
		// global cooldown. A total failure returns an error to trigger it.
		if len(errs) > 0 && successfulProjects == 0 {
			return errors.Join(errs...)
		}

		return nil
	}
}

// buildSupaBaseTask checks the Supabase projects and restores the ones that
// are inactive. A single token manages all the projects.
func buildSupaBaseTask(token string, store *status.Store) worker.Task {
	return func(ctx context.Context) error {
		log := logger.NewLogger("SUPABASE")

		checker := supabase.NewChecker(supabase.NewClient(token), log)

		projects, err := checker.Check()
		var services []status.MonitoredService
		var wsStatus status.Status
		var msg string
		if err != nil {
			wsStatus = status.Error
			msg = err.Error()
		} else {
			for _, p := range projects {
				services = append(services, status.MonitoredService{
					Target: p.Ref,
					State:  p.Status,
				})
			}
			wsStatus = status.Operational
		}
		ws, _ := status.NewWorkerStatus("supabase", "Supabase", wsStatus, services, msg)
		if store != nil {
			store.Set(ws)
		}
		return err
	}
}

func buildAternosTask(botService service.BotService, store *status.Store) worker.Task {
	return func(ctx context.Context) error {
		return botService.CheckAndStartServerWithStore(ctx, store)
	}
}
