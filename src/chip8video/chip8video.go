package chip8video

const HEIGTH = 32
const WIDTH = 64

type Video struct {
	pixelbuffer [HEIGTH][WIDTH]bool
}

// create new video driver with emtpy buffer
func CreateVideo() *Video {
	return new(Video)
}

// clear the screen
func Clear(video *Video) {
	//TBI
}