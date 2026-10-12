package aternos

import (
	"context"

	config "github.com/sebx771/monitor-bot/internal/config"
	"github.com/sebx771/monitor-bot/internal/integrations/aternos/adapters"
	"github.com/sebx771/monitor-bot/internal/integrations/aternos/minecraft"
	service "github.com/sebx771/monitor-bot/internal/integrations/aternos/services"
	"github.com/sebx771/monitor-bot/internal/status"
	"github.com/sebx771/monitor-bot/internal/worker"
)

// NewTask encapsulates the Aternos/Minecraft wiring and returns the worker task.
func NewTask(at *config.AternosConfig, store *status.Store) worker.Task {
	checker := minecraft.NewChecker(at.GetHost(), at.GetPort())
	gistClient := adapters.NewGitHubGistClient(at.GetGithubToken(), at.GetGistID())
	botService := service.NewBotService(
		checker,
		gistClient,
		at.GetStoragePath(),
		at.GetServerID(),
		at.IsHeadless(),
		at.GetHost(),
	)

	return func(ctx context.Context) error {
		return botService.CheckAndStartServerWithStore(ctx, store)
	}
}
