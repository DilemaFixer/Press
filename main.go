package main

import (
	win "github.com/DilemaFixer/Press/src/window"
	"time"
)

func main(){
	window := win.NewWindow(800,600,60,"Hello World!")
	time.Sleep(10*time.Second)
	window.Destroy()
}
