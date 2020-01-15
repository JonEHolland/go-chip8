package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	var state = newState("blinky.ch8")
	var timers = newTimers()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Go-Chip8",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		64,
		32,
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

	//sdl.PollEvent()
	//sdl.Delay(2000)

	window.UpdateSurface()

	// TODO - This loop should run at a specific speed
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		executeCycle(state, timers)
		timers.executeTimers()

		if state.drawFlag {
			renderer.Clear()
			for x := uint8(0); x < 64; x++ {
				for y := uint8(0); y < 32; y++ {
					if state.graphicsBuffer[x+(x*y)] == 1 {
						renderer.DrawPoint(int32(x), int32(y))
					}
				}
			}

			renderer.Present()
			window.UpdateSurface()
			state.drawFlag = false
		}

		// TODO - Set Key State
	}
}
