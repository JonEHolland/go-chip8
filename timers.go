package main

// Timers struct
type Timers struct {
	counter    uint8
	delayTimer uint8
	soundTimer uint8
}

func newTimers() *Timers {
	return &Timers{}
}

func (t *Timers) executeTimers() {
	// The CPU loops runs at 500hz, and the timers run at 60hz.
	// Return early if it has not been enough CPU cycles to update the timers yet
	// 500hz is 1 cycle every 2ms, 60hz is 1 cycle every 16ms, so use 8 as the counter to
	// get close.
	if t.counter < 8 {
		t.counter++
		return
	}

	t.counter = 0

	if t.delayTimer > 0 {
		t.delayTimer--
	}

	if t.soundTimer > 0 {
		// A Beep should play
		t.soundTimer--
	}
}
