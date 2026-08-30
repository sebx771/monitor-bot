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

	"github.com/sebx771/monitor-bot/internal/aiven"
	config "github.com/sebx771/monitor-bot/internal/config"
	"github.com/sebx771/monitor-bot/internal/logger"

	"github.com/sebx771/monitor-bot/internal/adapters"
	"github.com/sebx771/monitor-bot/internal/supabase"
	"github.com/sebx771/monitor-bot/internal/minecraft"
	service "github.com/sebx771/monitor-bot/internal/services"



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

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("error loading configuration", "error", err)
		os.Exit(1)
	}

	// Listen for Ctrl+C to stop workers cleanly
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	aivenEnabled := cfg.GetAivenConfig().IsEnabled()
	aternosEnabled := cfg.GetAternosConfig().IsEnabled()
	supabaseEnabled:=  cfg.GetSupaBaseConfig().IsEnabled()

	log.Info("active modules", "|AIVEN|", aivenEnabled, "|ATERNOS|", aternosEnabled , "|SUPABASE|", supabaseEnabled)

	if !aivenEnabled && !aternosEnabled && !supabaseEnabled {
		log.Warn("no module enabled, no workers to run")
		return
	}

	var wg sync.WaitGroup

	// Aternos (Minecraft) bot — only if enabled
	if aternosEnabled {
		at := cfg.GetAternosConfig()

		checker := minecraft.NewChecker(at.GetHost(), at.GetPort())
		gistClient := adapters.NewGitHubGistClient(at.GetGithubToken(), at.GetGistID())
		botService := service.NewBotService(checker, gistClient,
			at.GetStoragePath(), at.GetServerID(), at.IsHeadless())

		w, err := worker.New(checkInterval, errCooldown, botService.CheckAndStartServer)
		if err != nil {
			log.Error("error initializing Aternos worker", "error", err)
			os.Exit(1)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := w.Run(ctx); err != nil {
				log.Error("Minecraft worker finished with error", "error", err)
			}
		}()
	}

	// Aiven monitor — only if enabled
	if aivenEnabled {
		aivenTask := buildAivenTask(cfg.GetAivenConfig().GetCredentials())

		wAiven, err := worker.New(aivenInterval, aivenCooldown, aivenTask)
		if err != nil {
			log.Error("error initializing Aiven worker", "error", err)
			os.Exit(1)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := wAiven.Run(ctx); err != nil {
				log.Error("Aiven worker finished with error", "error", err)
			}
		}()
	}

	if supabaseEnabled {
		supabaseTask := buildSupaBaseTask(cfg.GetSupaBaseConfig().GetToken())

		wSupaBase, err := worker.New(aivenInterval, aivenCooldown, supabaseTask)
		if err != nil {
			log.Error("error initializing Supabase worker", "error", err)
			os.Exit(1)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := wSupaBase.Run(ctx); err != nil {
				log.Error("Supabase worker finished with error", "error", err)
			}
		}()
	}

	wg.Wait()

	log.Info("application finished successfully")
}

// buildAivenTask runs a checker per credential. An API failure does not abort
// the others; if at least one fails, it returns a combined error to trigger
// the worker's cooldown.
func buildAivenTask(credentials []config.Credential) worker.Task {
	return func(ctx context.Context) error {
		var errs []error

		for _, cred := range credentials {
			checker := aiven.NewChecker(aiven.NewClient(cred.Token), cred.Project)

			if err := checker.Check(); err != nil {
				errs = append(errs, fmt.Errorf("project %s: %w", cred.Project, err))
				continue
			}
		}

		return errors.Join(errs...)
	}
}

// buildSupaBaseTask checks the Supabase projects and restores the ones that
// are inactive. A single token manages all the projects.
func buildSupaBaseTask(token string) worker.Task {
	return func(ctx context.Context) error {
		log := logger.NewLogger("SUPABASE")

		checker := supabase.NewChecker(supabase.NewClient(token), log)

		return checker.Check()
	}
}
