package main

import (
	"chip8cpu"
	"chip8mem"
	"fmt"
	"time"
)

func main() {
	fmt.Println("[>] Starting emulator")
	fmt.Println("[>] Loading ROM")
	cpu := chip8cpu.CreateCpu()
	err := chip8mem.LoadROM(cpu.Mem, "BC_test.ch8")
	if err != nil {
		fmt.Println("[!] Error when loading ROM: ", err)
	}

	for true {
		err := chip8cpu.Tick(cpu)
		if err != nil {
			fmt.Println("[!] CPU has thrown an error: ", err)
			break
		}
		time.Sleep(16666 * time.Microsecond) // T=1/60=16 2/3 ms for 60Hz
	}

	fmt.Println("[>] Emulator done, good bye")
}
