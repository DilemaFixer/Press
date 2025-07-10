package window

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
    "sync"
)

type Window struct {
    sdlWindow   *sdl.Window
    renderer    *sdl.Renderer
    width       int32
    height      int32
    fps         int32
    isDestroyed bool
}

var (
    sdlInitialized bool
    sdlMutex       sync.Mutex
)

func InitSDL() error {
    sdlMutex.Lock()
    defer sdlMutex.Unlock()
    
    if sdlInitialized {
        return nil 
    }
    if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
        return fmt.Errorf("failed to initialize SDL: %w", err)
    }
    
    sdlInitialized = true
    return nil
}

func NewWindow(width, height, fps int32, title string) (*Window, error) {
    if !sdlInitialized {
        return nil, fmt.Errorf("SDL not initialized, call InitSDL() first")
    }
    
    sdlWindow, err := sdl.CreateWindow(
        title,
        sdl.WINDOWPOS_UNDEFINED,
        sdl.WINDOWPOS_UNDEFINED,
        width,
        height,
        sdl.WINDOW_SHOWN,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create window: %w", err)
    }
    
    renderer, err := sdl.CreateRenderer(
        sdlWindow,
        -1,
        sdl.RENDERER_ACCELERATED,
    )
    if err != nil {
        sdlWindow.Destroy()  
		return nil, fmt.Errorf("failed to create renderer: %w", err)
    }
    
    return &Window{
        sdlWindow: sdlWindow,
        renderer:  renderer,
        width:     width,
        height:    height,
        fps:       fps,
    }, nil
}

func (w *Window) Destroy() {
    if w.isDestroyed {
        return 
    }
    
    if w.renderer != nil {
        w.renderer.Destroy()
    }
    if w.sdlWindow != nil {
        w.sdlWindow.Destroy()
    }
    
    w.isDestroyed = true
}

func (w *Window) GetRenderer() *sdl.Renderer {
    return w.renderer
}

func (w *Window) GetSDLWindow() *sdl.Window {
    return w.sdlWindow
}

func (w *Window) GetSize() (int32, int32) {
    return w.width, w.height
}

func (w *Window) GetFPS() int32 {
    return w.fps
}

func QuitSDL() {
    sdlMutex.Lock()
    defer sdlMutex.Unlock()
    
    if sdlInitialized {
        sdl.Quit()
        sdlInitialized = false
    }
}
