package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/mxschmitt/playwright-go"
)

type Browser struct {
	pw  *playwright.Playwright
	browser  playwright.Browser
	context playwright.BrowserContext

	headless bool
	running bool
}

func NewBrowser(headless bool) *Browser {
	return &Browser{
		headless: headless,
		running: false,
	}
}


func (b *Browser) Start(ctx context.Context) error {
	if b.running  {
		return nil
	}
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("error iniciando Playwright %w: ", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(b.headless),
	})

	if err != nil {
		pw.Stop()
		return fmt.Errorf("error iniciando Chromium: %w", err)
	}



	b.pw = pw
	b.browser = browser
	b.running = true

	return nil
}

func (b *Browser) LoadStorageState(path string) error {
	if !b.running {
		return errors.New("browser no iniciado")
	}

	// Si existe un contexto anterior, lo cerramos
	if b.context != nil {
		if err := b.context.Close(); err != nil {
			return fmt.Errorf("cerrando contexto anterior: %w", err)
		}

		b.context = nil
	}

	// Crear nuevo contexto con el storage state
	context, err := b.browser.NewContext(
		playwright.BrowserNewContextOptions{
			StorageStatePath: playwright.String(path),
		},
	)

	if err != nil {
		return fmt.Errorf("creando contexto con storage state: %w", err)
	}

	b.context = context

	return nil
}


func (b *Browser) SaveStorageState(path string) error {
	if !b.running {
		return errors.New("browser no iniciado")
	}

	if b.context == nil {
		return errors.New("no existe un contexto activo")
	}

	_,err := b.context.StorageState(
		playwright.BrowserContextStorageStateOptions{
			Path: playwright.String(path),
		},
	)

	if err != nil {
		return fmt.Errorf("guardando storage state: %w", err)
	}

	return nil
}


func (b *Browser) NewPage() (playwright.Page, error) {
    if !b.running {
        return nil, errors.New("browser no iniciado")
    }

    return b.context.NewPage()
}


func (b *Browser) Stop() error {
	// cerramos de forma jerarquica los componentes playwright
	var stopErr error

	if err := b.context.Close(); err != nil {
		stopErr = fmt.Errorf("cerrando contexto: %w", err)
	}

	if err := b.browser.Close(); err != nil && stopErr == nil {
		stopErr = fmt.Errorf("cerrando browser: %w", err)
	}

	if err := b.pw.Stop(); err != nil && stopErr == nil {
		stopErr = fmt.Errorf("deteniendo playwright: %w", err)
	}

	b.running = false
	b.context = nil
	b.browser = nil
	b.pw = nil

	return stopErr
}

func (b *Browser) Restart(ctx context.Context) error{
	if !b.running{
		return fmt.Errorf("El Browser se encuentra apagado")
	}
    err:= b.Stop()
	if err != nil {
		return fmt.Errorf("error al apagar Browser: %w",err)
	}
	err = b.Start(ctx)
	if err != nil{
		return fmt.Errorf("error al iniciar Browser: %w ",err)
	}

	return nil
}



func (b *Browser) IsRunning() bool {
	return b.running
}
