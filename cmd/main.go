package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sebx771/monitor-bot/internal/adapters"
	"github.com/sebx771/monitor-bot/internal/automation"
	"github.com/sebx771/monitor-bot/internal/minecraft"
)

const (
	onlinePollInterval = 15 * time.Second
	onlineWaitTimeout  = 5 * time.Minute
)

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

func waitForOnline(checker *minecraft.Checker) error {
	deadline := time.Now().Add(onlineWaitTimeout)

	for time.Now().Before(deadline) {
		online, err := checker.IsOnline()
		if err == nil && online {
			return nil
		}

		log.Printf("servidor aún offline, reintentando en %s...", onlinePollInterval)
		time.Sleep(onlinePollInterval)
	}

	return fmt.Errorf("el servidor no se puso online en %s", onlineWaitTimeout)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := loadEnv(".env"); err != nil {
		return err
	}

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		return err
	}

	checker := minecraft.NewChecker(os.Getenv("HOST"), uint16(port))

	online, err := checker.IsOnline()
	if err != nil {
		log.Println("Servidor no responde al ping:", err)
		online = false
	}

	if online {
		log.Println("el servidor ya está online")
		return nil
	}

	storagePath := os.Getenv("STORAGE_PATH")

	if err := ensureStorageFile(storagePath); err != nil {
		return err
	}

	ctx := context.Background()

	browser := adapters.NewBrowser()

	if err := browser.Start(ctx); err != nil {
		return err
	}

	if err := browser.LoadStorageState(storagePath); err != nil {
		return err
	}

	defer func() {
		if err := browser.Stop(); err != nil {
			log.Printf("error cerrando browser: %v", err)
		}
	}()

	defer func() {
		if err := browser.SaveStorageState(storagePath); err != nil {
			log.Printf("error guardando sesión: %v", err)
		}
	}()

	bot := automation.NewAternosBot(browser)

	if err := bot.Open(); err != nil {
		return err
	}

	logged, err := bot.IsLogged()
	if err != nil {
		return err
	}

	if !logged {
		return fmt.Errorf("sesión de Aternos no válida: inicia sesión manualmente en el navegador abierto y vuelve a ejecutar el bot")
	}

	if err := bot.StartServer(os.Getenv("SERVER_ID")); err != nil {
		return err
	}

	log.Println("click de inicio enviado, esperando a que el servidor esté online...")

	if err := waitForOnline(checker); err != nil {
		return err
	}

	log.Println("el servidor está online")

	return nil
}
