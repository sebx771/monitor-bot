package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sebx771/monitor-bot/internal/adapters"
	"github.com/sebx771/monitor-bot/internal/automation"
	"github.com/sebx771/monitor-bot/internal/minecraft"
)

type BotService struct {
	checker     *minecraft.Checker
	storagePath string
	serverID    string
}

func NewBotService(checker *minecraft.Checker, storagePath, serverID string) *BotService {
	return &BotService{
		checker:     checker,
		storagePath: storagePath,
		serverID:    serverID,
	}
}

// CheckAndStartServer es la función que cumple con el tipo `Task` del Worker
func (s *BotService) CheckAndStartServer(ctx context.Context) error {
	log.Println("[Service] Verificando estado del servidor Minecraft...")

	online, err := s.checker.IsOnline()
	if err != nil {
		log.Println("[Service] Servidor no responde al ping:", err)
		online = false
	}

	if online {
		log.Println("[Service] El servidor ya está online. Omitiendo encendido.")
		return nil
	}

	log.Println("[Service] Servidor offline. Iniciando automatización de Aternos...")

	// Solo levantamos el navegador si el servidor REALMENTE está offline
	if err := s.startAternosServer(ctx); err != nil {
		return fmt.Errorf("falló el encendido del servidor: %w", err)
	}

	log.Println("[Service] Clic de inicio enviado, esperando a que el servidor esté online...")

	if err := s.waitForOnline(); err != nil {
		return fmt.Errorf("el servidor no confirmó estar online: %w", err)
	}

	log.Println("[Service] ¡El servidor ya se encuentra online y listo para jugar!")
	return nil
}

// startAternosServer encapsula la apertura y cierre del navegador de forma segura
func (s *BotService) startAternosServer(ctx context.Context) error {
	browser := adapters.NewBrowser()

	if err := browser.Start(ctx); err != nil {
		return err
	}
	defer func() {
		if err := browser.Stop(); err != nil {
			log.Printf("error cerrando browser: %v", err)
		}
	}()

	if err := browser.LoadStorageState(s.storagePath); err != nil {
		return err
	}

	defer func() {
		if err := browser.SaveStorageState(s.storagePath); err != nil {
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
		return fmt.Errorf("sesión de Aternos no válida: inicia sesión manualmente")
	}

	return bot.StartServer(s.serverID)
}

func (s *BotService) waitForOnline() error {
	onlinePollInterval := 15 * time.Second
	onlineWaitTimeout := 5 * time.Minute

	deadline := time.Now().Add(onlineWaitTimeout)

	for time.Now().Before(deadline) {
		online, err := s.checker.IsOnline()
		if err == nil && online {
			return nil
		}
		log.Printf("[Service] Servidor aún offline, reintentando en %s...", onlinePollInterval)
		time.Sleep(onlinePollInterval)
	}

	return fmt.Errorf("el servidor no se puso online en %s", onlineWaitTimeout)
}