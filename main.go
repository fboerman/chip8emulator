package main

import (
	"chip8cpu"
	"chip8keyboard"
	"chip8mem"
	"chip8video"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	ROM_fname := flag.String("ROM", "", "name of ROM file to load and start")
	debug := flag.Bool("debug", false, "dump PC and instr in hex format for each cycle")
	ticktime := flag.Int("ticktime", 50, "time in ms between ticks")

	flag.Parse()

	if _, err := os.Stat(*ROM_fname); os.IsNotExist(err) {
		fmt.Println("[!] Invalid ROM file!")
		return
	}

	fmt.Println("[>] Starting emulator")
	fmt.Println("[>] Loading ROM")
	cpu := chip8cpu.CreateCpu()
	defer chip8video.CloseVideo(cpu.Video)

	err := chip8mem.LoadROM(cpu.Mem, *ROM_fname)
	if err != nil {
		fmt.Println("[!] Error when loading ROM: ", err)
	}
	chip8mem.LoadFonts(cpu.Mem)

	fmt.Println("[>] Running video test")
	// run test by first displaying F manually and then loading B from font
	chip8video.Test(cpu.Video, 1, []uint8{})
	chip8video.Render(cpu.Video)
	time.Sleep(1 * time.Second)
	sprite := chip8mem.LoadFontSprite(cpu.Mem, 0xB)
	chip8video.Test(cpu.Video, 2, sprite)
	chip8video.Render(cpu.Video)
	time.Sleep(1 * time.Second)
	chip8video.Clear(cpu.Video)
	fmt.Println("[>] Video test done")

	fmt.Println("[>] Starting CPU loop")
	for true {
		if *debug {
			chip8cpu.DebugDump(cpu)
		}
		err := chip8cpu.Tick(cpu)
		if err != nil {
			fmt.Print("[!] CPU has thrown an error: ", err)
			fmt.Printf(" at PC 0x%X\n", cpu.Mem.PC)
			break
		}
		if cpu.Video.Dirty {
			chip8video.Render(cpu.Video)
			cpu.Video.Dirty = false
		}
		// process the keyboard
		chip8keyboard.Update(cpu.Keyboard)
		// time between ticks
		time.Sleep(time.Duration(*ticktime) * time.Millisecond)
		//time.Sleep(16666 * time.Microsecond) // T=1/60=16 2/3 ms for 60Hz
	}
	fmt.Println("[>] Emulator done, good bye")
}
