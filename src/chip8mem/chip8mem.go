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
const FONTSTART = 0x50 // start address for fonts

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

// load n seperate bytes from memory
func LoadnBytes(mem *Memory, addr uint16, n int) (data []uint8) {
	for i := 0; i < n; i++ {
		data = append(data, mem.mem[addr+uint16(i)])
	}

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

// load the standard sprites for fonts
func LoadFonts(mem *Memory) {
	var fontset = [...]uint8{
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

	for i := range fontset {
		// do not user writebyte as we are writing to the part of memory that normal programs are not allowed to load to
		mem.mem[uint16(FONTSTART+i)] = fontset[i]
	}
}
