package chip8cpu

import (
	"chip8mem"
	"chip8video"
	"errors"
	"fmt"
	"math"
	"math/rand"
)

type Cpu struct {
	Mem   *chip8mem.Memory
	Video *chip8video.Video
}

// create new CPU, emtpy initialized
func CreateCpu() *Cpu {
	cpu := new(Cpu)
	cpu.Mem = chip8mem.CreateMem()
	cpu.Video = chip8video.CreateVideo()

	return cpu
}

//execute instruction from current PC
// variables used in comments about the instruction:
//nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
//n or nibble - A 4-bit value, the lowest 4 bits of the instruction
//x - A 4-bit value, the lower 4 bits of the high byte of the instruction
//y - A 4-bit value, the upper 4 bits of the low byte of the instruction
//n - A 4-bit value, the lower 4 bits of the low byte of the instruction
//kk or byte - An 8-bit value, the lowest 8 bits of the instruction
func Tick(cpu *Cpu) error {
	// load current instruction and extract its upper 4 bits as opcode
	var instr uint16 = chip8mem.LoadInstr(cpu.Mem, cpu.Mem.PC)
	var opcode uint8 = uint8(instr >> 12)

	//select 12 lower bits (0xFFF = 0000 1111 1111 1111)
	var nnn uint16 = instr & 0xFFF
	// lower 4 bits of high byte (0xF = 1111)
	var x uint8 = uint8(instr>>8) & 0xF
	// upper 4 bits of low byte
	y := uint8((instr >> 4) & 0xF)
	// get value encoded in instruction in lower byte so bitmask with ( 0xFF = 1111 1111)
	var kk uint8 = uint8(instr & 0xFF)
	// lower 4 bits of the low byte of the instruction
	var n uint8 = uint8(instr & 0xF)

	// act on opcode
	switch opcode {
	case 0:
		//get function code and act on it
		//bitmask only last 12 bits
		functioncode := nnn
		switch functioncode {
		case 0xE0:
			// CLS
			// Clear the display

			chip8video.Clear(cpu.Video)
			cpu.Mem.PC += 2
		case 0xEE:
			// RET
			// Return from a subroutine

			new_addr, err := chip8mem.PopStack(cpu.Mem)
			if err != nil {
				return err
			}
			cpu.Mem.PC = new_addr
		default:
			return errors.New(fmt.Sprintf("Malformed instruction (0x%X), wrong functioncode (0x%X) with opcode (0x%X)", instr, functioncode, opcode)) // TODO: make this custom error type
		}
	case 1:
		// JP addr
		// Jump to location nnn

		cpu.Mem.PC = nnn
	case 2:
		// CALL addr
		// Call subroutine at nnn

		err := chip8mem.AddStack(cpu.Mem, cpu.Mem.PC)
		if err != nil {
			return err
		}
		cpu.Mem.PC = nnn
	case 3:
		// SE Vx, byte
		// Skip next instruction if Vx = kk

		v, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}

		if *v == kk {
			cpu.Mem.PC += 4
		} else {
			cpu.Mem.PC += 2
		}
	case 4:
		// SNE Vx, byte
		// exact opposite of SE: Skip next instruction if Vx != kk

		v, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}

		if *v != kk {
			cpu.Mem.PC += 4
		} else {
			cpu.Mem.PC += 2
		}
	case 5:
		// SE Vx, Vy
		// Skip next instruction if Vx = Vy

		Vx, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}
		Vy, err := chip8mem.GetReg(cpu.Mem, y)
		if err != nil {
			return err
		}
		if *Vx == *Vy {
			cpu.Mem.PC += 4
		} else {
			cpu.Mem.PC += 2
		}
	case 6:
		// LD Vx, byte
		// Set Vx = kk

		Vx, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}
		// write low byte from instruction
		*Vx = kk
		cpu.Mem.PC += 2
	case 7:
		// ADD Vx, byte
		//Set Vx = Vx + kk

		Vx, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}
		// write low byte from instruction
		*Vx += kk
		cpu.Mem.PC += 2
	case 8:

		Vx, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}
		Vy, err := chip8mem.GetReg(cpu.Mem, y)
		if err != nil {
			return err
		}

		// function code in the last 4 bits so bitmask with (0xF == 1111)
		functioncode := uint8(instr & 0xF)
		VF, _ := chip8mem.GetReg(cpu.Mem, 0xF)
		switch functioncode {
		case 0:
			// LD Vx, Vy
			// Set Vx = Vy
			*Vx = *Vy
		case 1:
			// OR Vx, Vy
			//Set Vx = Vx OR Vy
			*Vx = *Vx | *Vy
		case 2:
			// AND Vx, Vy
			// Set Vx = Vx AND Vy
			*Vx = *Vx & *Vy
		case 3:
			// XOR Vx, Vy
			// Set Vx = Vx XOR Vy
			*Vx = *Vx ^ *Vy
		case 4:
			// ADD Vx, Vy
			// Set Vx = Vx + Vy, set VF = carry

			temp := uint16(*Vx) + uint16(*Vy)
			if temp > math.MaxUint8 {
				*VF = 1
			} else {
				*VF = 0
			}
			// write the lowest byte of the result
			*Vx = uint8(temp & 0xFF)
		case 5:
			// SUB Vx, Vy
			// Set Vx = Vx - Vy, set VF = NOT borrow

			if *Vx > *Vy {
				*VF = 1
			} else {
				*VF = 0
			}

			*Vx = *Vx - *Vy
		case 6:
			// SHR Vx {, Vy}
			// Set Vx = Vx SHR 1

			*VF = *Vx & 0x1
			*Vx = *Vx >> 1
		case 7:
			// SUBN Vx, Vy
			// Set Vx = Vy - Vx, set VF = NOT borrow

			if *Vy > *Vx {
				*VF = 1
			} else {
				*VF = 0
			}
			*Vx = *Vy - *Vx
		case 0xE:
			// SHL Vx {, Vy}
			// Set Vx = Vx SHL 1

			*VF = *Vx & 0x80
			*Vx = *Vx << 1

		default:
			return errors.New(fmt.Sprintf("Malformed instruction (0x%X), wrong functioncode (0x%X) with opcode (0x%X)", instr, functioncode, opcode)) // TODO: make this custom error type
		}

		cpu.Mem.PC += 2

	case 9:
		// SNE Vx, Vy
		// Skip next instruction if Vx != Vy

		Vx, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}
		Vy, err := chip8mem.GetReg(cpu.Mem, y)
		if err != nil {
			return err
		}

		if *Vx != *Vy {
			cpu.Mem.PC += 4
		} else {
			cpu.Mem.PC += 2
		}

	case 0xA:
		//  LD I, addr
		// Set I = nnn
		cpu.Mem.I = nnn
		cpu.Mem.PC += 2
	case 0xB:
		// JP V0, addr
		// Jump to location nnn + V0
		V0, _ := chip8mem.GetReg(cpu.Mem, 0x0)
		cpu.Mem.PC = (instr & 0xFFF) + uint16(*V0)
	case 0xC:
		// RND Vx, byte
		// Set Vx = random byte AND kk

		Vx, err := chip8mem.GetReg(cpu.Mem, x)

		if err != nil {
			return err
		}

		random := uint8(rand.Intn(math.MaxUint8 + 1))
		*Vx = random & uint8(instr&0xFF)
		cpu.Mem.PC += 2
	case 0xD:
		// DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision

		Vx, err := chip8mem.GetReg(cpu.Mem, x)
		if err != nil {
			return err
		}
		Vy, err := chip8mem.GetReg(cpu.Mem, y)
		if err != nil {
			return err
		}
		sprite := chip8mem.LoadnBytes(cpu.Mem, cpu.Mem.I, int(n))
		VF, _ := chip8mem.GetReg(cpu.Mem, 0xF)
		*VF = chip8video.DisplaySprite(cpu.Video, sprite, *Vx, *Vy)

		cpu.Mem.PC += 2
	case 0xE:
		// SKP Vx
		// Skip next instruction if key with the value of Vx is not pressed

		// TBI
		return errors.New(fmt.Sprintf("Instruction 0x(%X) has non implemented opcode 0x(%X)", instr, opcode))
	case 0xF:
		functioncode := kk

		Vx, err := chip8mem.GetReg(cpu.Mem, x)

		if err != nil {
			return err
		}

		switch functioncode {
		case 7:
			// LD Vx, DT
			// Set Vx = delay timer value
			*Vx = cpu.Mem.T_delay
		case 0xA:
			// LD Vx, K
			// Wait for a key press, store the value of the key in Vx

			// TBI
			return errors.New(fmt.Sprintf("Instruction 0x(%X) has non implemented opcode 0x(%X) and functioncode 0x(%X)", instr, opcode, functioncode))
		case 0x15:
			// LD DT, Vx
			// Set delay timer = Vx

			cpu.Mem.T_delay = *Vx
		case 0x18:
			// LD ST, Vx
			// Set sound timer = Vx

			cpu.Mem.T_sound = *Vx
		case 0x1E:
			// ADD I, Vx
			// Set I = I + Vx

			cpu.Mem.I += uint16(*Vx)
		case 0x29:
			// LD F, Vx
			// Set I = location of sprite for digit Vx

			cpu.Mem.I = uint16(chip8mem.FONTSTART + *Vx*5)

			cpu.Mem.PC += 2
		case 0x33:
			// LD B, Vx
			// Store BCD representation of Vx in memory locations I, I+1, and I+2
			// this formula is ugly in that it is hard to replicate in hardware
			// TODO: rewrite this into something that is closer to hardware units

			temp := *Vx
			// ones-place
			err := chip8mem.WriteByte(cpu.Mem, cpu.Mem.I+2, temp%10)
			if err != nil {
				return err
			}
			temp /= 10

			// tens-place
			err = chip8mem.WriteByte(cpu.Mem, cpu.Mem.I+1, temp%10)
			if err != nil {
				return err
			}
			temp /= 10

			// hundreds-place
			err = chip8mem.WriteByte(cpu.Mem, cpu.Mem.I, temp%10)
			if err != nil {
				return err
			}

		case 0x55:
			// LD [I], Vx
			// Store registers V0 through Vx in memory starting at location I
			if x >= chip8mem.NUMREGS {
				return errors.New(fmt.Sprintf("Invalid reg number %d", x))
			}

			for i := 0; i < int(x)+1; i++ {
				V, _ := chip8mem.GetReg(cpu.Mem, uint8(i))
				err := chip8mem.WriteByte(cpu.Mem, cpu.Mem.I+uint16(i), *V)
				if err != nil {
					return err
				}
			}

		case 0x65:
			// LD Vx, [I]
			// Read registers V0 through Vx from memory starting at location I
			if x >= chip8mem.NUMREGS {
				return errors.New(fmt.Sprintf("Invalid reg number %d", x))
			}
			if cpu.Mem.I+uint16(x) > chip8mem.MEMSIZE {
				return errors.New(fmt.Sprintf("Invalid address 0x(%X) to read to memory", cpu.Mem.I+uint16(x)))
			}
			if cpu.Mem.I < chip8mem.MEMSTART {
				return errors.New(fmt.Sprintf("Invalid address 0x(%X) to read to memory", cpu.Mem.I))
			}

			for i := 0; i < int(x)+1; i++ {
				V, _ := chip8mem.GetReg(cpu.Mem, uint8(i))
				*V = chip8mem.LoadByte(cpu.Mem, cpu.Mem.I+uint16(i))
			}
		}
		cpu.Mem.PC += 2
	default:
		return errors.New(fmt.Sprintf("Instruction 0x(%X), with non recognized opcode 0x(%X)", instr, opcode))
	}

	return nil
}
