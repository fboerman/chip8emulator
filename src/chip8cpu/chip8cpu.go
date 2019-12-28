package chip8cpu

import (
	"chip8mem"
	"chip8video"
	"errors"
	"fmt"
)

type Cpu struct {
	Mem *chip8mem.Memory
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
func Tick(cpu *Cpu) error {
	// load current instruction and extract its upper 4 bits as opcode
	var instr uint16 = chip8mem.LoadInstr(cpu.Mem, cpu.Mem.PC)
	var opcode uint8 = uint8(instr >> 12)

	// act on opcode
	switch opcode {
	case 0:
		//get function code and act on it
		//bitmask only last 12 bits
		functioncode := instr & 0xFFF
		switch functioncode {
		case 0xE0:
			// CLS
			chip8video.Clear(cpu.Video)
			cpu.Mem.PC += 2
		case 0xEE:
			// RET
			new_addr, err := chip8mem.PopStack(cpu.Mem)
			if err != nil {
				return err
			}
			cpu.Mem.PC = new_addr
		default:
			return errors.New(fmt.Sprintf("Malformed instruction (0x%X), wrong functioncode (0x%X) with opcode (0x%X)", instr, functioncode, opcode)) // TODO: make this custom error type
		}
	case 1:
		// JP
		//select 12 lower bits (0xFFF = 0000 1111 1111 1111)
		cpu.Mem.PC = instr & 0xFFF
	case 2:
		// CALL


	default:
		return errors.New(fmt.Sprintf("Instruction (%X), with non recognized opcode (%X)", instr, opcode))
	}

	return nil
}