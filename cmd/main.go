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
	"github.com/sebx771/monitor-bot/internal/minecraft"
	service "github.com/sebx771/monitor-bot/internal/services"



	"github.com/sebx771/monitor-bot/internal/worker"
)

const (
	checkInterval = 1440 * time.Minute // Frecuencia de revisión del servidor
	errCooldown   = 70 * time.Minute   // Tiempo de espera si falla Aternos

	aivenInterval = 60 * time.Minute // Frecuencia de revisión de servicios Aiven
	aivenCooldown = 30 * time.Minute // Tiempo de espera si falla la API de Aiven
)

func main() {
	log := logger.NewLogger("MAIN")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("error cargando configuración", "error", err)
		os.Exit(1)
	}

	// Escuchar Ctrl+C para detener los workers limpiamente
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	aivenEnabled := cfg.GetAivenConfig().IsEnabled()
	aternosEnabled := cfg.GetAternosConfig().IsEnabled()

	log.Info("módulos activos|", "aiven", aivenEnabled, "aternos", aternosEnabled)

	if !aivenEnabled && !aternosEnabled {
		log.Warn("ningún módulo habilitado, no hay workers que ejecutar")
		return
	}

	var wg sync.WaitGroup

	// Bot de Aternos (Minecraft) — solo si está habilitado
	if aternosEnabled {
		at := cfg.GetAternosConfig()

		checker := minecraft.NewChecker(at.GetHost(), at.GetPort())
		gistClient := adapters.NewGitHubGistClient(at.GetGithubToken(), at.GetGistID())
		botService := service.NewBotService(checker, gistClient,
			at.GetStoragePath(), at.GetServerID(), at.IsHeadless())

		w, err := worker.New(checkInterval, errCooldown, botService.CheckAndStartServer)
		if err != nil {
			log.Error("error inicializando worker de Aternos", "error", err)
			os.Exit(1)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := w.Run(ctx); err != nil {
				log.Error("worker de Minecraft finalizó con error", "error", err)
			}
		}()
	}

	// Monitor de Aiven — solo si está habilitado
	if aivenEnabled {
		aivenTask := buildAivenTask(cfg.GetAivenConfig().GetCredentials())

		wAiven, err := worker.New(aivenInterval, aivenCooldown, aivenTask)
		if err != nil {
			log.Error("error inicializando worker de Aiven", "error", err)
			os.Exit(1)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := wAiven.Run(ctx); err != nil {
				log.Error("worker de aiven finalizó con error", "error", err)
			}
		}()
	}

	wg.Wait()

	log.Info("aplicación finalizada correctamente")
}

// buildAivenTask ejecuta un checker por credencial. Un fallo en una API no
// aborta a las demás; si al menos una falla, retorna un error combinado para
// activar el cooldown del worker.
func buildAivenTask(credentials []config.Credential) worker.Task {
	return func(ctx context.Context) error {
		var errs []error

		for _, cred := range credentials {
			checker := aiven.NewChecker(aiven.NewClient(cred.Token), cred.Project)

			if err := checker.Check(); err != nil {
				errs = append(errs, fmt.Errorf("proyecto %s: %w", cred.Project, err))
				continue
			}
		}

		return errors.Join(errs...)
	}
}
