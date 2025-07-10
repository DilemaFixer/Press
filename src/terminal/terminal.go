package terminal

import (
	"fmt"
	"log"
	"time"

	"github.com/DilemaFixer/Press/src/window"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Terminal struct {
	window     *window.Window
	font       *ttf.Font
	width      int32
	height     int32
	charWidth  int32
	charHeight int32
	rows       int32
	cols       int32

	buffer      [][]rune
	cursorX     int32
	cursorY     int32
	running     bool
	inputBuffer string
}

type Color struct {
	R, G, B, A uint8
}

var (
	ColorBlack   = Color{0, 0, 0, 255}
	ColorWhite   = Color{255, 255, 255, 255}
	ColorGreen   = Color{0, 255, 0, 255}
	ColorRed     = Color{255, 0, 0, 255}
	ColorBlue    = Color{0, 0, 255, 255}
	ColorYellow  = Color{255, 255, 0, 255}
	ColorCyan    = Color{0, 255, 255, 255}
	ColorMagenta = Color{255, 0, 255, 255}
)

func NewTerminal(width, height, fps int32, pathToFont string) (*Terminal, error) {
	if err := window.InitSDL(); err != nil {
		return nil, fmt.Errorf("failed to initialize SDL: %w", err)
	}

	if err := ttf.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize TTF: %w", err)
	}

	win, err := window.NewWindow(width, height, fps, "Terminal")
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	font, err := ttf.OpenFont(pathToFont, 16)
	if err != nil {
		win.Destroy()
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	charWidth, charHeight, err := font.SizeUTF8("W")
	if err != nil {
		font.Close()
		win.Destroy()
		return nil, fmt.Errorf("failed to get font size: %w", err)
	}

	rows := height / int32(charHeight)
	cols := width / int32(charWidth)

	terminal := &Terminal{
		window:      win,
		font:        font,
		height:      height,
		width:       width,
		charWidth:   int32(charWidth),
		charHeight:  int32(charHeight),
		cols:        cols,
		rows:        rows,
		buffer:      make([][]rune, rows),
		cursorX:     0,
		cursorY:     0,
		running:     true, 
		inputBuffer: "",
	}

	for i := range terminal.buffer {
		terminal.buffer[i] = make([]rune, cols)
		for j := range terminal.buffer[i] {
			terminal.buffer[i][j] = '*' 
		}
	}

	return terminal, nil
}

func (t *Terminal) Render() error {
	renderer := t.window.GetRenderer()

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	for row := int32(0); row < t.rows; row++ {
		for col := int32(0); col < t.cols; col++ {
			char := t.buffer[row][col]
			if char == ' ' {
				continue
			}

			if err := t.renderChar(char, col, row, ColorWhite); err != nil {
				return err
			}
		}
	}

	t.renderCursor()

	renderer.Present()
	return nil
}

func (t *Terminal) renderChar(char rune, x, y int32, color Color) error {
	if char == ' ' {
		return nil
	}

	renderer := t.window.GetRenderer()

	surface, err := t.font.RenderUTF8Solid(string(char), sdl.Color{
		R: color.R,
		G: color.G,
		B: color.B,
		A: color.A,
	})
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	rect := &sdl.Rect{
		X: x * t.charWidth,
		Y: y * t.charHeight,
		W: t.charWidth,
		H: t.charHeight,
	}

	return renderer.Copy(texture, nil, rect)
}

func (t *Terminal) renderCursor() {
	renderer := t.window.GetRenderer()

	if time.Now().UnixMilli()%1000 < 500 {
		renderer.SetDrawColor(255, 255, 255, 255)
		rect := &sdl.Rect{
			X: t.cursorX * t.charWidth,
			Y: t.cursorY * t.charHeight,
			W: t.charWidth,
			H: t.charHeight,
		}
		renderer.FillRect(rect)
	}
}

func (t *Terminal) IsRunning() bool {
	return t.running
}

func (t *Terminal) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			t.running = false
			log.Println("Quit event received")

		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					t.running = false
					log.Println("ESC pressed")
				case sdl.K_RETURN:
					t.handleEnter()
				case sdl.K_BACKSPACE:
					t.handleBackspace()
				}
			}

		case *sdl.TextInputEvent:
			t.handleTextInput(e.GetText())
		}
	}
}

func (t *Terminal) handleEnter() {
	t.cursorX = 0
	t.cursorY++
	if t.cursorY >= t.rows {
		t.cursorY = t.rows - 1
	}
}

func (t *Terminal) handleBackspace() {
	if t.cursorX > 0 {
		t.cursorX--
		t.buffer[t.cursorY][t.cursorX] = ' '
	}
}

func (t *Terminal) handleTextInput(text string) {
	for _, char := range text {
		if t.cursorX < t.cols {
			t.buffer[t.cursorY][t.cursorX] = char
			t.cursorX++
		}
	}
}

func (t *Terminal) WriteText(text string, color Color) {
	for _, char := range text {
		if char == '\n' {
			t.handleEnter()
			continue
		}

		if t.cursorX >= t.cols {
			t.handleEnter()
		}

		if t.cursorY >= t.rows {
			break
		}

		t.buffer[t.cursorY][t.cursorX] = char
		t.cursorX++
	}
}

func (t *Terminal) Clear() {
	for i := range t.buffer {
		for j := range t.buffer[i] {
			t.buffer[i][j] = ' '
		}
	}
	t.cursorX = 0
	t.cursorY = 0
}

func (t *Terminal) Destroy() {
	if t.font != nil {
		t.font.Close()
	}
	if t.window != nil {
		t.window.Destroy()
	}
}
