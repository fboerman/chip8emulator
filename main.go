package main

import (
	"chip8cpu"
	"chip8mem"
	"chip8video"
	"fmt"
	"time"
)

func main() {
	fmt.Println("[>] Starting emulator")
	fmt.Println("[>] Loading ROM")
	cpu := chip8cpu.CreateCpu()
	defer chip8video.CloseVideo(cpu.Video)

	err := chip8mem.LoadROM(cpu.Mem, "BC_test.ch8")
	if err != nil {
		fmt.Println("[!] Error when loading ROM: ", err)
	}

	fmt.Println("[>] Running video test")
	chip8video.Test(cpu.Video)
	chip8video.Render(cpu.Video)
	time.Sleep(5 * time.Second)

	//fmt.Println("[>] Starting CPU loop")
	//for true {
	//	err := chip8cpu.Tick(cpu)
	//	if err != nil {
	//		fmt.Print("[!] CPU has thrown an error: ", err)
	//		fmt.Printf(" at PC 0x%X\n", cpu.Mem.PC)
	//		break
	//	}
	//	if cpu.Video.Dirty {
	//		chip8video.Render(cpu.Video)
	//		cpu.Video.Dirty = false
	//	}
	//	time.Sleep(16666 * time.Microsecond) // T=1/60=16 2/3 ms for 60Hz
	//}

	fmt.Println("[>] Emulator done, good bye")
}
