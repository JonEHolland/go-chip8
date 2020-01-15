package main

func main() {
	var state = newState("blinky.ch8")
	var timers = newTimers()

	// TODO - This loop should run at a specific speed
	for {
		executeCycle(state, timers)
		timers.executeTimers()

		if state.drawFlag {
			// TODO Update SDL Buffer
		}

		// TODO - Set Key State
	}
}
