package chip8mem

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
)

// some constants from definition
const MEMSIZE = 4096   // main memory size
const NUMREGS = 16     // how many general purpose 8-bit registers are available
const STACKSIZE = 16   // ammount of 16-bits registers the stack is composed of
const MEMSTART = 0x200 // start address of all memory

type Memory struct {
	mem     [MEMSIZE]uint8
	regs    [NUMREGS]uint8
	stack   [STACKSIZE]uint16
	PC      uint16
	SP      uint8
	T_delay uint8  // delay timer
	T_sound uint8  // sound timer
	I       uint16 // index register
}

// initialize empty memory
func CreateMem() *Memory {
	mem := new(Memory)
	mem.PC = MEMSTART
	mem.SP = math.MaxUint8 // start stack pointer at underflow, then if increased by 1 it points to element 0
	return mem
}

// load rom from file into memory, overwrite what was there already
func LoadROM(mem *Memory, fname string) error {
	// open the file
	file, err := os.Open(fname)

	if err != nil {
		return err
	}
	defer file.Close()

	//read at most the available memory space of bytes
	bytes := make([]byte, MEMSIZE-MEMSTART)
	reader := bufio.NewReader(file)
	_, err = reader.Read(bytes)

	if err != nil {
		return err
	}

	//copy into the memory struct
	copy(mem.mem[MEMSTART:], bytes)

	return nil
}

// load 2 seperate bytes from memory
func Load2Bytes(mem *Memory, addr uint16) (data [2]uint8) {
	data[0] = mem.mem[addr]
	data[1] = mem.mem[addr+1]

	return
}

// load 2 bytes from memory concatenated
func LoadInstr(mem *Memory, addr uint16) uint16 {
	return (uint16(mem.mem[addr]) << 8) | uint16(mem.mem[addr+1])
}

// load 1 byte from memory
func LoadByte(mem *Memory, addr uint16) uint8 {
	return mem.mem[addr]
}

// pop address from the stack and adjust the stack pointer
// note that this does not actually clear the stack register, only the stackpointer
// return stackunderflow error if stack is empty
func PopStack(mem *Memory) (addr uint16, err error) {
	if mem.SP == math.MaxUint8 {
		return 0, errors.New("Stack underflow")
	}
	addr = mem.stack[mem.SP]
	mem.SP--
	return
}

// add address to the stack and adjust stack pointer
// return stackoverflow error if stack is full
func AddStack(mem *Memory, addr uint16) (err error) {
	if mem.SP == STACKSIZE-1 {
		// SP already points to top of stack
		return errors.New("Stack overflow")
	}
	mem.SP++
	mem.stack[mem.SP] = addr

	return
}

// get pointer to register
// return error if invalid register number
func GetReg(mem *Memory, x uint8) (v *uint8, err error) {
	if x >= NUMREGS {
		return nil, errors.New(fmt.Sprintf("Invalid reg number %d", x))
	}
	v = &(mem.regs[x])
	return

}

//write byte to address in memory
func WriteByte(mem *Memory, addr uint16, byte uint8) error {
	if addr >= MEMSIZE || addr < MEMSTART {
		return errors.New(fmt.Sprintf("Invalid address 0x(%X) to write to memory", addr))
	}

	mem.mem[addr] = byte

	return nil
}
