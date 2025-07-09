package window

import (
	"github.com/veandco/go-sdl2/sdl"

	"log"
)

type Window struct {
	window *sdl.Window
	renderer *sdl.Renderer

	width int32
	height int32
	fps int32
}

func NewWindow(width ,height int32 , fps int32 ,title string) *Window {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Failed to initialize SDL: %v", err)
	}
	
	window, err := sdl.CreateWindow(
		title,           
		sdl.WINDOWPOS_UNDEFINED,        
		sdl.WINDOWPOS_UNDEFINED,        
		width,                    
		height,                   	
		sdl.WINDOW_SHOWN,               
	)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	
	renderer, err := sdl.CreateRenderer(
		window,                         
		-1,                            
		sdl.RENDERER_ACCELERATED,      
	)

	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}

	return &Window{
		window:window,
		renderer:renderer,
		height:height,
		width:width,
		fps:fps,
	}
}

func (window *Window)Destroy(){
	window.renderer.Destroy() 
	window.window.Destroy() 
	sdl.Quit() 
}

