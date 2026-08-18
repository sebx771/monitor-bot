package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sebx771/monitor-bot/internal"
	"github.com/sebx771/monitor-bot/internal/adapters"
	"github.com/sebx771/monitor-bot/internal/aiven"
	"github.com/sebx771/monitor-bot/internal/minecraft"
	service "github.com/sebx771/monitor-bot/internal/services"
	"github.com/sebx771/monitor-bot/internal/worker"
)

const (
	checkInterval = 1440 * time.Minute // Frecuencia de revisión del servidor 
	errCooldown   = 70 * time.Minute // Tiempo de espera si falla Aternos

	aivenInterval = 60 * time.Minute // Frecuencia de revisión de servicios Aiven
	aivenCooldown = 30 * time.Minute // Tiempo de espera si falla la API de Aiven
)

func main() {
	cfg, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Escuchar Ctrl+C para detener los workers limpiamente
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

    //configuracion del bot para aternos 
	checker := minecraft.NewChecker(cfg.GetHostMC(), cfg.GetPortMC())
	gistClient := adapters.NewGitHubGistClient(cfg.GetGithubToken(), cfg.GetGistID())
	botService := service.NewBotService(checker, gistClient, cfg.GetStoragePath(), cfg.GetServerID(), cfg.GetHeadless())

	w, err := worker.New(checkInterval, errCooldown, botService.CheckAndStartServer)
	if err != nil {
		log.Fatalf("Error al inicializar el worker: %v", err)
	}
    

	aivenTask := buildAivenTask(cfg.GetAivenCredentials())
	wAiven, err := worker.New(aivenInterval, aivenCooldown, aivenTask)
	if err != nil {
		log.Fatalf("Error al inicializar el worker de Aiven: %v", err)
	}

	log.Printf("Iniciando Monitor Bot (Intervalo: %s, Cooldown: %s)...", checkInterval, errCooldown)
	log.Printf("Iniciando Worker Aiven (Intervalo: %s, Cooldown: %s)...", aivenInterval, aivenCooldown)

	//  Activación de los Workers en paralelo
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := w.Run(ctx); err != nil {
			log.Printf("Worker de Minecraft finalizó con error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := wAiven.Run(ctx); err != nil {
			log.Printf("Worker de Aiven finalizó con error: %v", err)
		}
	}()

	wg.Wait()

	log.Println("Aplicación finalizada correctamente.")
}

// buildAivenTask ejecuta un checker por credencial. Un fallo en una API no
// aborta a las demás; si al menos una falla, retorna un error combinado para
// activar el cooldown del worker.
func buildAivenTask(credentials []internal.Credential) worker.Task {
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
