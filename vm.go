package main

import (
	"fmt"
	"math/rand"
)

func executeCycle(s *State, t *Timers) {
	// Opcodes are two sequential bytes, so combine them into the currentOpcode
	s.currentOpcode = uint16(s.memory[s.programCounter])<<8 | uint16(s.memory[s.programCounter+1])
	var rX = (s.currentOpcode & 0x0F00) >> 8
	var rY = (s.currentOpcode & 0x00F0) >> 4

	fmt.Printf("%.4X - ", s.programCounter)

	switch s.currentOpcode & 0xF000 {

	case 0x0000:

		switch s.currentOpcode {

		case 0x00E0: // 00EE - Clear Screen
			fmt.Println("CLS")
			for i := 0; i < 2048; i++ {
				s.graphicsBuffer[i] = 0
			}
			s.drawFlag = true
			break

		case 0x00EE: // 00EE - Return from Subroutine
			fmt.Println("RET")
			s.stackPointer--
			s.programCounter = s.stack[s.stackPointer]
			break

		default:
			fmt.Printf("0x0000 - Unimplemented Opcode: %.4X\n", s.currentOpcode)
		}
		break

	case 0x1000: // 1NNN - Jump to address NNN
		fmt.Printf("JP %.4X\n", s.currentOpcode&0x0FFF)
		s.programCounter = (s.currentOpcode & 0x0FFF)
		return

	case 0x2000: // 2NNN - Call Subroutine at NNN
		fmt.Printf("CALL %.4X\n", s.currentOpcode&0x0FFF)
		s.stack[s.stackPointer] = s.programCounter
		s.stackPointer++
		s.programCounter = (s.currentOpcode & 0x0FFF)
		return

	case 0x3000: // 3XNN - Skip next instruction if RegisterX == NN
		fmt.Printf("SE %d %.4X\n", rX, (s.currentOpcode & 0x00FF))
		var value = (s.currentOpcode & 0x00FF)
		if s.registers[rX] == uint8(value) {
			s.programCounter += 2
		}
		break

	case 0x4000: // 3XNN - Skip next instruction if RegisterX != NN
		fmt.Printf("SNE %d %.4X\n", rX, (s.currentOpcode & 0x00FF))
		var value = (s.currentOpcode & 0x00FF)
		if s.registers[rX] != uint8(value) {
			s.programCounter += 2
		}
		break

	case 0x5000: // 5XY0 - Skip next instruction if RegisterX == RegisterY
		fmt.Printf("SE %d %d\n", rX, rY)
		if s.registers[rX] == s.registers[rY] {
			s.programCounter += 2
		}
		break

	case 0x6000: // 6XNN - Sets Register X to NN
		fmt.Printf("LD %d %.4X\n", rX, s.currentOpcode&0x00FF)
		s.registers[rX] = uint8(s.currentOpcode & 0x00FF)
		break

	case 0x7000: // 7XNN - Adds NN to Register X
		fmt.Printf("ADD %d %.4X\n", rX, s.currentOpcode&0x00FF)
		s.registers[rX] += uint8(s.currentOpcode & 0x00FF)
		break

	case 0x8000:
		switch s.currentOpcode & 0x000F {

		case 0x0000: // 8XY0 - Sets Register X to the value of Register Y
			fmt.Printf("LD %d %d\n", rX, rY)
			s.registers[rX] = s.registers[rY]
			break
		case 0x0001: // 8XY1 - Sets Register X to Register X | Register Y
			fmt.Printf("OR %d %d\n", rX, rY)
			s.registers[rX] |= s.registers[rY]
			break
		case 0x0002: // 8XY1 - Sets Register X to Register X & Register Y
			fmt.Printf("AND %d %d\n", rX, rY)
			s.registers[rX] &= s.registers[rY]
			break
		case 0x0003: // 8XY2 - Sets Register X to RegisterX ^ Register Y
			fmt.Printf("XOR %d %d\n", rX, rY)
			s.registers[rX] ^= s.registers[rY]
			break
		case 0x0004: // 8XY4 - Add Register X to Register Y, setting Register F if there is a carry
			fmt.Printf("ADD %d %d\n", rX, rY)
			var temp = uint16(s.registers[rX] + s.registers[rY])
			if temp > 255 {
				s.registers[0xF] = 1
			} else {
				s.registers[0xF] = 0
			}

			if temp > 255 {
				temp -= 256
			}
			s.registers[rX] = uint8(temp)
			break
		case 0x0005: // 8XY5 - Sub Register X from Register Y, unsetting Register F if there is a borrow
			fmt.Printf("SUB %d %d\n", rX, rY)
			var temp = uint16(s.registers[rX] - s.registers[rY])
			if s.registers[rX] > s.registers[rY] {
				s.registers[0xF] = 1
			} else {
				s.registers[0xF] = 0
			}

			if temp < 0 {
				temp += 256
			}
			s.registers[rX] = uint8(temp)
			break
		case 0x0006: // 8XY6 - Shifts Register X right by 1, Resister F is set to the LSB before the shift
			fmt.Printf("SHR %d\n", rX)
			s.registers[0xF] = s.registers[rX] & 0x1
			s.registers[rX] >>= 1
			break
		case 0x0007: // 8XY7 - Sub Register Y from Register X - Register X. Register F is unset if there is a borrow
			fmt.Printf("SUBN %d %d\n", rX, rY)
			var temp = uint16(s.registers[rY] - s.registers[rX])
			if s.registers[rY] > s.registers[rX] {
				s.registers[0xF] = 1
			} else {
				s.registers[0xF] = 0
			}

			if temp < 0 {
				temp += 256
			}
			s.registers[rX] = uint8(temp)
			break
		case 0x000E: // 8XYE - Shifts Register X left by 1, Resister F is set to the MSB before the shift
			fmt.Printf("SHL %d\n", rX)
			s.registers[0xF] = s.registers[rX] >> 7
			s.registers[rX] <<= 1
			break
		default:
			fmt.Printf("0x8000 - Unimplemented Opcode: %.4X\n", s.currentOpcode)
			break
		}
		break

	case 0x9000: // 5XY0 - Skip next instruction if RegisterX != RegisterY
		fmt.Printf("SNE %d %d\n", rX, rY)
		if s.registers[rX] != s.registers[rY] {
			s.programCounter += 2
		}
		break

	case 0xA000: // ANNN - Sets index pointer to the address NNN
		fmt.Printf("LD I %.4X\n", (s.currentOpcode&0x0FFF)>>8)
		s.indexPointer = (s.currentOpcode & 0x0FFF) >> 8
		break

	case 0xB000: // BNNN - Jumps to the address at NNN + Register 0
		fmt.Printf("JP 0 %.4X\n", s.currentOpcode&0x0FFF)
		s.programCounter = ((s.currentOpcode & 0x0FFF) + uint16(s.registers[0]))
		return

	case 0xC000: // CXNN - Sets Register X to a Random Number & NN
		fmt.Printf("RND %d %.4X\n", rX, s.currentOpcode&0x0FFF)
		var mask = (s.currentOpcode & 0x00FF) >> 8
		s.registers[rX] = uint8(rand.Intn(255)) & uint8(mask)
		break

	case 0xD000: // DXYN - Draws a Sprite
		// TODO - Implement Graphics
		fmt.Printf("DRW\n")
		s.drawFlag = true
		break

	case 0xE000:
		switch s.currentOpcode & 0x00FF {

		case 0x009E: // EX9E - Skips next instruction if key stored in Register X is pressed
			fmt.Printf("SKP %d\n", rX)
			if s.keyStates[s.registers[rX]] != 0 {
				s.programCounter += 2
			}
			break
		case 0x00A1: // EXA1 - Skips next instruction if the key stored in Register X isn't pressed
			fmt.Printf("SKPN %d\n", rX)
			if s.keyStates[s.registers[rX]] == 0 {
				s.programCounter += 2
			}
			break
		default:
			fmt.Printf("0xE000 - Unimplemented Opcode: %.4X\n", s.currentOpcode)
			break
		}

	case 0xF000:
		{
			switch s.currentOpcode & 0x00FF {

			case 0x0007: // FX07 - Sets Register X to the delay timer
				fmt.Printf("LD %d DT\n", rX)
				s.registers[rX] = t.delayTimer
				break
			case 0x000A: // FX0A - Wait for a key press, then store it in Register X
				fmt.Printf("LD %d K\n", rX)
				var keyPressed = false
				for i := 0; i < 16; i++ {
					if s.keyStates[i] != 0 {
						s.registers[rX] = uint8(i)
						keyPressed = true
					}
				}

				// Only increment the program counter if a key was found, this leaves
				// execution blocked on this opcode until a key is pressed
				if !keyPressed {
					s.programCounter -= 2
				}
				break
			case 0x0015: // FX15 - Sets the delay timer to Register X
				fmt.Printf("LD DT %d\n", rX)
				t.delayTimer = s.registers[rX]
				break
			case 0x0018: // FX18 - Sets the sound timer to Register X
				fmt.Printf("LD ST %d\n", rX)
				t.soundTimer = s.registers[rX]
				break
			case 0x001E: // FX1E - Adds Register X to Index, sets Register F if overflows
				fmt.Printf("ADD I %d\n", rX)
				if s.indexPointer+uint16(s.registers[rX]) > 0xFFF {
					s.registers[0xF] = 1
				} else {
					s.registers[0xF] = 0
				}

				s.indexPointer += uint16(s.registers[rX])
				break
			case 0x0029: // FX29 - Sets Index to the sprite location for a character in Register X
				fmt.Printf("LD F %d\n", rX)
				s.indexPointer = uint16(s.registers[rX]) * 0x5
				break
			case 0x0033: // FX33 - Stores the BCD representation of Register X at I, I+1, I+2
				fmt.Printf("LD B %d\n", rX)
				var registerValue = s.registers[rX]
				s.memory[s.indexPointer] = registerValue / 100
				s.memory[s.indexPointer+1] = (registerValue / 10) % 10
				s.memory[s.indexPointer+2] = registerValue % 10
				break
			case 0x0055: // FX55 - Stores Register 0 to Register X at address I
				fmt.Printf("LD [i] %d\n", rX)
				var finalRegister = uint16((s.currentOpcode & 0x0F00) >> 8)

				for i := uint16(0); i <= finalRegister; i++ {
					s.memory[s.indexPointer+i] = s.registers[i]
				}

				// Move index pointer past all the stored registers
				s.indexPointer += finalRegister + 1
				break
			case 0x0065: // FX65 - Fills Register 0 to Register X with memory at address I
				fmt.Printf("LD %d [i]\n", rX)
				var finalRegister = uint16((s.currentOpcode & 0x0F00) >> 8)

				for i := uint16(0); i <= finalRegister; i++ {
					s.registers[i] = s.memory[s.indexPointer+i]
				}

				// Move index pointer past all the stored registers
				s.indexPointer += finalRegister + 1
				break
			default:
				fmt.Printf("0xF000 - Unimplemented Opcode: %.4X\n", s.currentOpcode)
			}
		}

	default:
		fmt.Printf("Default - Unimplemented Opcode: %.4X\n", s.currentOpcode)
		break
	}

	s.programCounter += 2
}
