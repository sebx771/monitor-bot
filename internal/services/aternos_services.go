package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sebx771/monitor-bot/internal/adapters"
	"github.com/sebx771/monitor-bot/internal/automation"
	"github.com/sebx771/monitor-bot/internal/logger"
	"github.com/sebx771/monitor-bot/internal/minecraft"
	port "github.com/sebx771/monitor-bot/internal/ports"
)

var log = logger.NewLogger("SERVICE")

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

// CheckAndStartServer is the function that satisfies the Worker's `Task` type
func (s *BotService) CheckAndStartServer(ctx context.Context) error {
	log.Info("checking Minecraft server status")

	online, err := s.checker.IsOnline()
	if err != nil {
		log.Warn("server does not respond to ping", "error", err)
		online = false
	}

	if online {
		log.Info("server is already online, skipping startup")
		return nil
	}

	log.Info("server offline, starting Aternos automation")

	// We only launch the browser if the server is REALLY offline
	if err := s.startAternosServer(ctx); err != nil {
		return fmt.Errorf("server startup failed: %w", err)
	}

	log.Info("start click sent, waiting for the server to come online")

	if err := s.waitForOnline(); err != nil {
		return fmt.Errorf("the server did not confirm it is online: %w", err)
	}

	log.Info("the server is now online and ready to play")
	return nil
}

// startAternosServer encapsulates opening and closing the browser safely
func (s *BotService) startAternosServer(ctx context.Context) error {
	// Best-effort: if the remote download fails, we continue with the local file
	_ = s.syncStorageFromRemote(ctx)

	browser := adapters.NewBrowser(s.headless)

	if err := browser.Start(ctx); err != nil {
		return err
	}
	defer func() {
		if err := browser.Stop(); err != nil {
			log.Error("error closing browser", "error", err)
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
		log.Warn("invalid Aternos session, retrying with fresh cookies from the Gist")

		if err := s.retryWithFreshState(ctx, browser, bot); err != nil {
			return err
		}
	}

	if err := bot.StartServer(s.serverID); err != nil {
		return err
	}

	// Only on the successful path: we save the session locally and sync it
	// with the Gist. It avoids overwriting good remote cookies with stale
	// ones from failed cycles.
	if err := browser.SaveStorageState(s.storagePath); err != nil {
		log.Error("error saving session", "error", err)
	}

	s.syncStorageToRemote(ctx)

	return nil
}

// syncStorageFromRemote downloads the session state from the remote storage
// (GitHub Gist) to the local file. The error is logged and propagated so the
// caller can decide: the current flow treats it as best-effort.
func (s *BotService) syncStorageFromRemote(ctx context.Context) error {
	if err := s.storage.DownloadState(ctx, s.storagePath); err != nil {
		log.Warn("could not download the remote state, using the local one", "error", err)
		return err
	}

	log.Info("session state downloaded from the Gist")
	return nil
}

// syncStorageToRemote uploads the local session state file to the remote
// storage (GitHub Gist). The error is logged and does not abort the cycle.
func (s *BotService) syncStorageToRemote(ctx context.Context) error {
	if err := s.storage.UploadState(ctx, s.storagePath); err != nil {
		log.Warn("could not upload the state to the Gist", "error", err)
		return err
	}

	log.Info("session state uploaded to the Gist")
	return nil
}

// retryWithFreshState re-downloads the session state from the Gist, reloads
// the browser context (which invalidates the current page) and checks the
// session again. Returns an error if it is still invalid.
func (s *BotService) retryWithFreshState(ctx context.Context, browser *adapters.Browser, bot *automation.AternosBot) error {
	if err := s.syncStorageFromRemote(ctx); err != nil {
		return fmt.Errorf("session retry failed: %w", err)
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
		return fmt.Errorf("invalid Aternos session: log in manually")
	}

	log.Info("Aternos session recovered from the Gist")
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
		log.Info("server still offline, retrying", "interval", onlinePollInterval)
		time.Sleep(onlinePollInterval)
	}

	return fmt.Errorf("the server did not come online in %s", onlineWaitTimeout)
}