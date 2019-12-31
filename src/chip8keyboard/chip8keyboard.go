package chip8keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Keyboard struct {
	keys_state [16]uint8
}

// initialize empty keyboard
func CreateKeyboard() *Keyboard {
	return new(Keyboard)
}

// keyboard mapping, return pointer to register for that key
func mapping(keyboard *Keyboard, event *sdl.KeyboardEvent) (reg *uint8, addr uint8) {
	switch sdl.GetScancodeName(event.Keysym.Scancode) {
	case "0":
		addr = 0x0
	case "1":
		addr = 0x1
	case "2":
		addr = 0x2
	case "3":
		addr = 0x3
	case "4":
		addr = 0x4
	case "5":
		addr = 0x5
	case "6":
		addr = 0x6
	case "7":
		addr = 0x7
	case "8":
		addr = 0x8
	case "9":
		addr = 0x9
	case "A":
		addr = 0xA
	case "B":
		addr = 0xB
	case "C":
		addr = 0xC
	case "D":
		addr = 0xD
	case "E":
		addr = 0xE
	case "F":
		addr = 0xF
	default:
		addr = math.MaxUint8
		return nil, addr
	}

	reg = &keyboard.keys_state[addr]

	return
}

// find either the next keyboard down event or empty the event queue
// update the state accordingly
func Update(keyboard *Keyboard) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		if event.GetType() == sdl.KEYDOWN || event.GetType() == sdl.KEYUP {
			// keyboard mapping
			reg, _ := mapping(keyboard, event.(*sdl.KeyboardEvent))
			if reg != nil {
				switch event.GetType() {
				case sdl.KEYDOWN:
					*reg = 1
				case sdl.KEYUP:
					*reg = 0
				}
			}
		}
	}
}

// return bool if specified key is pressed
func IsPressed(keyboard *Keyboard, key uint8) bool {
	return keyboard.keys_state[key] == 1
}

// wait for keypress, update keyboard state, return the pressed key
func WaitforKey(keyboard *Keyboard) (pressed uint8) {
	for true {
		event := sdl.WaitEvent()
		switch event.GetType() {
		case sdl.KEYDOWN:

		}
		if event.GetType() == sdl.KEYDOWN || event.GetType() == sdl.KEYUP {
			// keyboard mapping
			reg, addr := mapping(keyboard, event.(*sdl.KeyboardEvent))
			if reg != nil {
				switch event.GetType() {
				case sdl.KEYDOWN:
					*reg = 1
					return addr // only exit at key pres instead of release
				case sdl.KEYUP:
					*reg = 0
				}
			}
		}
	}
	return
}
