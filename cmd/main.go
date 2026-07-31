package main

import (
	"fmt"
	"github.com/mxschmitt/playwright-go"
	"log"
)

func main() {
	err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	})
	if err != nil {
		log.Fatalf("no se pudo instalar el driver de playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		panic(err)
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		StorageStatePath: playwright.String("storage/state.json"),
	})

	if err != nil {
		panic(err)
	}

	page, err := context.NewPage()
	if err != nil {
		panic(err)
	}
      
	// cambiamos el evento de confirmacion de carga(load) , por el evento domloaded
	_, err = page.Goto("https://aternos.org", playwright.PageGotoOptions{
	   WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	   Timeout: playwright.Float(60000),
	})
	if err != nil {
		panic(err)
	}

	page.WaitForTimeout(60000)
	_, err = context.StorageState(playwright.BrowserContextStorageStateOptions{
		Path: playwright.String("storage/state.json"),
	})

	if err != nil {
       panic(err)
	}

	fmt.Println("terminando ejecucion...")
}
