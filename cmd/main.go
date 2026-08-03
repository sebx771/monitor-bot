package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sebx771/monitor-bot/internal/minecraft"
	 service "github.com/sebx771/monitor-bot/internal/services"
	"github.com/sebx771/monitor-bot/internal/worker"
)

const (
	checkInterval = 1 * time.Minute // Frecuencia de revisión del servidor
	errCooldown   = 2 * time.Minute // Tiempo de espera si falla Aternos
)

func main() {
	if err := loadEnv(".env"); err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}

	// Escuchar Ctrl+C para detener el worker limpiamente
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		log.Fatalf("Puerto inválido: %v", err)
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if err := ensureStorageFile(storagePath); err != nil {
		log.Fatalf("Error preparando almacenamiento de sesión: %v", err)
	}

	
	checker := minecraft.NewChecker(os.Getenv("HOST"), uint16(port))
	botService := service.NewBotService(checker, storagePath, os.Getenv("SERVER_ID"))

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

// --- Funciones de soporte de lectura del sistema ---

func loadEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(value))
	}
	return nil
}

func ensureStorageFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte("{}"), 0o644)
}