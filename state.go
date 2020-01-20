package main

import (
	"fmt"
	"io/ioutil"
)

// State Struct
type State struct {
	registers      [16]uint8
	memory         [4096]uint8
	graphicsBuffer [64][32]uint8
	drawFlag       bool
	indexPointer   uint16
	stack          [16]uint16
	stackPointer   uint16
	keyStates      [16]uint8
	programCounter uint16
	currentOpcode  uint16
}

func newState(romName string) *State {

	var fontData = [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	s := State{
		programCounter: 0x200,
	}

	// Load Font Data
	for i := 0; i < 80; i++ {
		s.memory[i] = fontData[i]
	}

	// Load Rom
	data, err := ioutil.ReadFile(romName)
	if err != nil {
		panic("Failed to open ROM")
	}

	for i := 0; i < len(data); i++ {
		// Program data starts at 0x200 in memory
		s.memory[i+0x200] = data[i]
	}

	return &s
}

func (s *State) dump () {
	fmt.Println("Dumping VM State")
	fmt.Println("-----------------")
	fmt.Printf("Program Counter: %.4X\n", s.programCounter)
	fmt.Printf("Current Opcode: %.4X\n", s.currentOpcode)
	fmt.Println("-----------------")
	fmt.Printf("Index Pointer: %.4X\n", s.indexPointer)
	fmt.Printf("Stack Pointer: %d\n", s.stackPointer)
	fmt.Println("-----------------")
	for r := 0; r < 16; r++ {
		fmt.Printf("Register %.1X: %.4X   ---   Stack %.1X: %.4X\n", r, s.registers[r], r, s.stack[r])
	}
	fmt.Println("-----------------")
}
