package main

import (
	"log"
	"time"

	"github.com/DilemaFixer/Press/src/terminal"
	"github.com/DilemaFixer/Press/src/window"
)

func main() {
	term, err := terminal.NewTerminal(800, 600, 60, "/System/Library/Fonts/Menlo.ttc")
	if err != nil {
		log.Fatal(err)
	}
	defer term.Destroy()
	defer window.QuitSDL()

	term.Clear()
	term.WriteText("Hello, Terminal!\n", terminal.ColorGreen)
	term.WriteText("Press ESC to exit\n", terminal.ColorYellow)

	log.Println("Terminal created, entering main loop...")

	for term.IsRunning() {
		term.HandleEvents()

		if err := term.Render(); err != nil {
			log.Printf("Render error: %v", err)
		}

		time.Sleep(16 * time.Millisecond) // ~60 FPS
	}

	log.Println("Terminal closed")
}
