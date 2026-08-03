package automation

import (
	"fmt"
	"time"

	"github.com/mxschmitt/playwright-go"
	port "github.com/sebx771/monitor-bot/internal/ports"
)


func dialogSelectors() []string {
	return []string{
		"text=Yes, I accept the EULA.",
		"text=Okay",
		"text=Accept",
		"text=I agree",
	}
}

type AternosBot struct{
	browser port.BrowserManager
	page playwright.Page
}

func (a *AternosBot) GetUrl() string{
	return "https://aternos.org"
}

func NewAternosBot(browser port.BrowserManager) *AternosBot {
	return &AternosBot{
		browser: browser,
	}
}


func (a *AternosBot) Open() error {

	page, err := a.browser.NewPage()
	if err != nil {
		return err
	}

	_, err = page.Goto(
		a.GetUrl(),
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			Timeout:   playwright.Float(30000),
		},
	)

	if err != nil {
		return err
	}

	a.page = page

	return nil
}

func (a *AternosBot) IsLogged() (bool, error) {

	count, err := a.page.
		Locator("text=Login").
		Count()

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (a *AternosBot) OpenServers() error {

	return a.page.
		Locator(`nav .mod-signup`).
		Click()
}

func (a *AternosBot) SelectServer(serverID string) error {
    selector := fmt.Sprintf(`div.server-body[data-id="%s"]`, serverID)

    fmt.Println("[Bot] Haciendo clic en el servidor de la lista...")
    if err := a.page.Locator(selector).Click(); err != nil {
        return err
    }

    // ESPERA CLAVE: Esperamos a que la URL cambie a la página del panel (/server/)
    fmt.Println("[Bot] Esperando a que cargue la página del panel...")
    err := a.page.WaitForURL("**/server/**", playwright.PageWaitForURLOptions{
        Timeout: playwright.Float(15000), // 15 segundos máximo
    })
    if err != nil {
        fmt.Println("[Bot] Advertencia: La URL no cambió a tiempo, continuando...")
    }

    return nil
}
func (b *AternosBot) ClickStart() error {
    
    start := b.page.Locator("#start")

    fmt.Println("[Bot] Esperando a que el botón Start sea visible...")
    
   
    if err := start.WaitFor(playwright.LocatorWaitForOptions{
        State:   playwright.WaitForSelectorStateVisible,
        Timeout: playwright.Float(15000),
    }); err != nil {
        return fmt.Errorf("el botón #start nunca apareció en la pantalla: %w", err)
    }

    fmt.Println("[Bot] Botón Start localizado. Enviando clic...")

    // Usamos Force: true por si hay un banner transparente de publicidad sobre el botón
    if err := start.Click(playwright.LocatorClickOptions{
        Force: playwright.Bool(true),
    }); err != nil {
        return fmt.Errorf("error al hacer clic en start: %w", err)
    }

    fmt.Println("[Bot] Clic enviado exitosamente. Revisando si aparece modal de confirmación...")
    
 
    b.HandleDialogs()

    return nil
}

func (a *AternosBot) HandleDialogs() error {
    for _, selector := range dialogSelectors() {
        locator := a.page.Locator(selector)

        // SOLUCIÓN: Agregamos un Timeout de 3 segundos. 
        // Si el modal no aparece rápido, saltamos al siguiente en lugar de esperar 30s.
        if err := locator.WaitFor(
            playwright.LocatorWaitForOptions{
                State:   playwright.WaitForSelectorStateVisible,
                Timeout: playwright.Float(3000), // 3000 milisegundos = 3 segundos
            },
        ); err != nil {
            continue
        }

        if err := locator.Click(); err != nil {
            fmt.Println("No se pudo clickear el modal:", err)
        }

        time.Sleep(500 * time.Millisecond)
    }

    return nil
}

func (a *AternosBot) StartServer(serverID string) error {
    if err := a.Open(); err != nil {
        return err
    }

    if err := a.OpenServers(); err != nil {
        return err
    }

    if err := a.SelectServer(serverID); err != nil {
        return err
    }

    if err := a.ClickStart(); err != nil {
        return err
    }



    return nil
}


