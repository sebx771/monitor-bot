package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sebx771/monitor-bot/internal"
	"github.com/sebx771/monitor-bot/internal/minecraft"
	service "github.com/sebx771/monitor-bot/internal/services"
	"github.com/sebx771/monitor-bot/internal/worker"
)

const (
	checkInterval = 60 * time.Minute // Frecuencia de revisión del servidor
	errCooldown   = 70 * time.Minute // Tiempo de espera si falla Aternos
)

func main() {
	cfg, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Escuchar Ctrl+C para detener el worker limpiamente
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	checker := minecraft.NewChecker(cfg.GetHostMC(), cfg.GetPortMC())
	botService := service.NewBotService(checker, cfg.GetStoragePath(), cfg.GetServerID())

	w, err := worker.New(checkInterval, errCooldown, botService.CheckAndStartServer)
	if err != nil {
		log.Fatalf("Error al inicializar el worker: %v", err)
	}

	log.Printf("Iniciando Monitor Bot (Intervalo: %s, Cooldown: %s)...", checkInterval, errCooldown)

	//  Activación del Worker
	if err := w.Run(ctx); err != nil {
		log.Fatalf("Error durante la ejecución del worker: %v", err)
	}

	log.Println("Aplicación finalizada correctamente.")
}
