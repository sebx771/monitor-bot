package main

import (
	"fmt"
	"github.com/mxschmitt/playwright-go"
)

func main() {
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

	_, err = page.Goto("https://aternos.org")
	if err != nil {
		panic(err)
	}

	page.WaitForTimeout(10000)
	_, err = context.StorageState(playwright.BrowserContextStorageStateOptions{
		Path: playwright.String("storage/state.json"),
	})

	if err != nil {
       panic(err)
	}

	fmt.Println("terminando ejecucion...")
}
