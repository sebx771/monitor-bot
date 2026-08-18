package port

import (
	"context"
	"github.com/mxschmitt/playwright-go"
)

type BrowserManager interface {
	Start(ctx context.Context) error

	SaveStorageState(path string) error
	LoadStorageState(path string) error

	NewPage() (playwright.Page, error)

	Stop() error
	Restart(ctx context.Context) error

	IsRunning() bool
}
