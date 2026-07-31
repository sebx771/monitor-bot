package port

import (
	"context"
	"github.com/mxschmitt/playwright-go"
)

type BrowserMonitor interface {
	Start(ctx context.Context) error
	Stop() error
	Restart(ctx context.Context) error

	NewPage() (playwright.Page, error)

	SaveStorageState(path string) error
	LoadStorageState(path string) error

	IsRunning() bool
}
