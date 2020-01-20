package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

func main() {
	var defaultWindowWidth = int32(640)
	var defaultWindowHeight = int32(320)
	var romName = os.Args[1]

	var state = newState(romName)
	// Hook so that state is dumped if the program panics
	defer func() {
		if err := recover(); err != nil {
			state.dump()
			panic (err)
		}
	}()

	var timers = newTimers()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	initAudio()
	sdl.PauseAudio(true)
	defer sdl.CloseAudio()

	window, err := sdl.CreateWindow(
		"Go-Chip8 - " + romName,
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
				updateKeys(state, event)
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
	renderer.SetDrawColor(0,0,0,255)
	renderer.Clear()
	renderer.SetDrawColor(255,255,255,255)
	for x := uint8(0); x < 64; x++ {
		for y := uint8(0); y < 32; y++ {
			if state.graphicsBuffer[x][y] == 1 {
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

func updateKeys(state *State, event sdl.Event) {
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		var direction = uint8(0)
		if t.Type == sdl.KEYDOWN {
			direction = uint8(1)
		} else {
			direction = uint8(0)
		}
		
		switch t.Keysym.Sym {
			case sdl.K_1: state.keyStates[0x1] = direction; break
			case sdl.K_2: state.keyStates[0x2] = direction; break
			case sdl.K_3: state.keyStates[0x3] = direction; break
			case sdl.K_4: state.keyStates[0xC] = direction; break
			case sdl.K_q: state.keyStates[0x4] = direction; break
			case sdl.K_w: state.keyStates[0x5] = direction; break
			case sdl.K_e: state.keyStates[0x6] = direction; break
			case sdl.K_r: state.keyStates[0xD] = direction; break
			case sdl.K_a: state.keyStates[0x7] = direction; break
			case sdl.K_s: state.keyStates[0x8] = direction; break
			case sdl.K_d: state.keyStates[0x9] = direction; break
			case sdl.K_f: state.keyStates[0xE] = direction; break
			case sdl.K_z: state.keyStates[0xA] = direction; break
			case sdl.K_x: state.keyStates[0x0] = direction; break
			case sdl.K_c: state.keyStates[0xB] = direction; break
			case sdl.K_v: state.keyStates[0xF] = direction; break
		}
	break
	}

}
