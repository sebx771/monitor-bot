package automation

import (
	"fmt"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/sebx771/monitor-bot/internal/logger"
	port "github.com/sebx771/monitor-bot/internal/ports"
)

var log = logger.NewLogger("ATERNOS")

func dialogSelectors() []string {
	return []string{
		"text=Yes, I accept the EULA.",
		"text=Okay",
		"text=Accept",
		"text=I agree",
	}
}

type AternosBot struct {
	browser port.BrowserManager
	page    playwright.Page
}

func (a *AternosBot) GetUrl() string {
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
	url := a.page.URL()
	title, err := a.page.Title()

	log.Debug("current URL", "url", url)
	if err != nil {
		log.Error("error getting title", "error", err)
	} else {
		log.Debug("title obtained", "title", title)
	}

    if title == "Just a moment..." {
    return fmt.Errorf("Cloudflare showed a verification page; the Aternos panel did not load")
}

	return a.page.
		Locator(`nav .mod-signup`).
		Click()
}

func (a *AternosBot) SelectServer(serverID string) error {
	selector := fmt.Sprintf(`div.server-body[data-id="%s"]`, serverID)

	log.Info("clicking the server in the list", "server_id", serverID)
	if err := a.page.Locator(selector).Click(); err != nil {
		return err
	}

	// KEY WAIT: We wait for the URL to change to the panel page (/server/)
	log.Info("waiting for the panel page to load")
	err := a.page.WaitForURL("**/server/**", playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(15000), // 15 seconds max
	})
	if err != nil {
		log.Warn("the URL did not change in time, continuing")
	}

	return nil
}
func (b *AternosBot) ClickStart() error {

	start := b.page.Locator("#start")

	log.Info("waiting for the Start button to be visible")

	if err := start.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(15000),
	}); err != nil {
		return fmt.Errorf("the #start button never appeared on screen: %w", err)
	}

	log.Info("Start button located, sending click")

	// We use Force: true in case there is a transparent ad banner over the button
	if err := start.Click(playwright.LocatorClickOptions{
		Force: playwright.Bool(true),
	}); err != nil {
		return fmt.Errorf("error clicking start: %w", err)
	}

	log.Info("click sent successfully, checking if a confirmation modal appears")

	b.HandleDialogs()

	return nil
}

func (a *AternosBot) HandleDialogs() error {
	for _, selector := range dialogSelectors() {
		locator := a.page.Locator(selector)

		// SOLUTION: We add a 3 second timeout.
		// If the modal does not appear quickly, we skip to the next one instead of waiting 30s.
		if err := locator.WaitFor(
			playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(3000), // 3000 milliseconds = 3 seconds
			},
		); err != nil {
			continue
		}

		if err := locator.Click(); err != nil {
			log.Error("could not click the modal", "modal", selector, "error", err)
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
