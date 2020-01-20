package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	var defaultWindowWidth = int32(640)
	var defaultWindowHeight = int32(320)

	var state = newState("roms/blinky.ch8")
	var timers = newTimers()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Go-Chip8",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		defaultWindowWidth,
		defaultWindowHeight,
		sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.Clear()
	renderer.SetDrawColor(0, 0, 0, 0)
	window.UpdateSurface()

	// TODO - This loop should run at a specific speed
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.KeyboardEvent:
				updateKeys(state, &event)
				break
			case *sdl.WindowEvent:
				// TODO - Handle Window Resize
				break
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		executeCycle(state, timers)
		timers.executeTimers()

		if state.drawFlag {
			drawScreen(state, window, renderer)
			state.drawFlag = false
		}
	}
}

func drawScreen(state *State, window *sdl.Window, renderer *sdl.Renderer) {
	renderer.Clear()
	renderer.SetDrawColor(0,0,0,255)
	for x := uint8(0); x < 64; x++ {
		for y := uint8(0); y < 32; y++ {
			if state.graphicsBuffer[x][y] == 1 {
				renderer.SetDrawColor(255,255,255,255)
				_ = renderer.FillRect(&sdl.Rect{
					X: int32(x) * 10,
					Y: int32(y) * 10,
					W: 10,
					H: 10})
			}
		}
	}

	renderer.Present()
	window.UpdateSurface()
}
func updateKeys(state *State, event *sdl.Event) {

}
