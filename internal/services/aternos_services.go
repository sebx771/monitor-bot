package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sebx771/monitor-bot/internal/adapters"
	"github.com/sebx771/monitor-bot/internal/automation"
	"github.com/sebx771/monitor-bot/internal/minecraft"
	port "github.com/sebx771/monitor-bot/internal/ports"
)

type BotService struct {
	checker     *minecraft.Checker
	storage     port.StateStorage
	storagePath string
	serverID    string
	headless    bool
}

func NewBotService(checker *minecraft.Checker, storage port.StateStorage, storagePath, serverID string, headless bool) *BotService {
	return &BotService{
		checker:     checker,
		storage:     storage,
		storagePath: storagePath,
		serverID:    serverID,
		headless:    headless,
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
	// Best-effort: si la descarga remota falla, se continúa con el archivo local
	_ = s.syncStorageFromRemote(ctx)

	browser := adapters.NewBrowser(s.headless)

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

	bot := automation.NewAternosBot(browser)

	if err := bot.Open(); err != nil {
		return err
	}

	logged, err := bot.IsLogged()
	if err != nil {
		return err
	}

	if !logged {
		log.Println("[Service] Sesión de Aternos no válida. Reintentando con cookies frescas del Gist...")

		if err := s.retryWithFreshState(ctx, browser, bot); err != nil {
			return err
		}
	}

	if err := bot.StartServer(s.serverID); err != nil {
		return err
	}

	// Solo en el camino exitoso: guardamos la sesión localmente y la
	// sincronizamos con el Gist. Evita pisar cookies remotas buenas con
	// cookies vencidas de ciclos fallidos.
	if err := browser.SaveStorageState(s.storagePath); err != nil {
		log.Printf("error guardando sesión: %v", err)
	}

	s.syncStorageToRemote(ctx)

	return nil
}

// syncStorageFromRemote descarga el estado de sesión desde el almacenamiento
// remoto (GitHub Gist) al archivo local. El error se loguea y se propaga para
// que el llamador decida: el flujo actual lo trata como best-effort.
func (s *BotService) syncStorageFromRemote(ctx context.Context) error {
	if err := s.storage.DownloadState(ctx, s.storagePath); err != nil {
		log.Printf("[Service] No se pudo descargar el estado remoto (se usa el local): %v", err)
		return err
	}

	log.Println("[Service] Estado de sesión descargado desde el Gist.")
	return nil
}

// syncStorageToRemote sube el archivo local de estado de sesión al
// almacenamiento remoto (GitHub Gist). El error se loguea y no aborta el ciclo.
func (s *BotService) syncStorageToRemote(ctx context.Context) error {
	if err := s.storage.UploadState(ctx, s.storagePath); err != nil {
		log.Printf("[Service] No se pudo subir el estado al Gist: %v", err)
		return err
	}

	log.Println("[Service] Estado de sesión subido al Gist.")
	return nil
}

// retryWithFreshState re-descarga el estado de sesión desde el Gist, recarga
// el contexto del navegador (lo que invalida la página actual) y vuelve a
// verificar la sesión. Retorna error si sigue siendo inválida.
func (s *BotService) retryWithFreshState(ctx context.Context, browser *adapters.Browser, bot *automation.AternosBot) error {
	if err := s.syncStorageFromRemote(ctx); err != nil {
		return fmt.Errorf("reintento de sesión fallido: %w", err)
	}

	if err := browser.LoadStorageState(s.storagePath); err != nil {
		return err
	}

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

	log.Println("[Service] Sesión de Aternos recuperada desde el Gist.")
	return nil
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