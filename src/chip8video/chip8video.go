package chip8video

import "github.com/veandco/go-sdl2/sdl"

const HEIGTH = 32
const WIDTH = 64
const SCALE = 10 // box of 10 pixels drawn with the actual pixel as relative 0,0 in top left

type Video struct {
	pixels   [HEIGTH][WIDTH]bool // direct pixels from program, false is white, true is black
	window   *sdl.Window
	renderer *sdl.Renderer
	tex      *sdl.Texture
	Dirty    bool
}

// create new video driver with emtpy buffer
func CreateVideo() *Video {
	video := new(Video)
	InitVideo(video)
	return video
}

// clear the buffer
func Clear(video *Video) {
	for y := 0; y < HEIGTH; y++ {
		for x := 0; x < WIDTH; x++ {
			video.pixels[y][x] = false
		}
	}
	video.renderer.SetDrawColor(0, 0, 0, 0)
	video.renderer.Clear()
}

// initialize SDL system and window
func InitVideo(video *Video) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("CHIP8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCALE*WIDTH, SCALE*HEIGTH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, int32(SCALE*WIDTH), int32(SCALE*HEIGTH))
	if err != nil {
		panic(err)
	}
	video.renderer.SetDrawColor(0, 0, 0, 0)
	video.window = window
	video.renderer = renderer
	video.tex = tex
}

// neatly close off SDL
func CloseVideo(video *Video) {
	video.tex.Destroy()
	video.renderer.Destroy()
	video.window.Destroy()
	sdl.Quit()
}

// render the current pixelbuffer to the texture
func Render(video *Video) {
	//updatebuffer(video)
	video.renderer.SetDrawColor(0, 0, 0, 0)
	video.renderer.Clear()
	//video.tex.Update(nil, video.pixelbuffer, WIDTH*SCALE*3)
	//video.renderer.Copy(video.tex, nil, nil)
	var rects []sdl.Rect
	for y := 0; y < HEIGTH; y++ {
		for x := 0; x < WIDTH; x++ {
			if video.pixels[y][x] {
				rects = append(rects, sdl.Rect{
					X: int32(x * SCALE),
					Y: int32(y * SCALE),
					W: SCALE,
					H: SCALE,
				})
			}
		}
	}
	video.renderer.SetDrawColor(255, 255, 255, 255)
	video.renderer.FillRects(rects)
	video.renderer.Present()
}

// render somethings on screen as test
func Test(video *Video) {
	video.pixels[0][0] = true
	video.pixels[0][1] = true
	video.pixels[0][2] = true
	video.pixels[0][3] = true
	video.pixels[1][0] = true
	video.pixels[2][0] = true
	video.pixels[2][1] = true
	video.pixels[2][2] = true
	video.pixels[3][0] = true
	video.pixels[4][0] = true
	video.pixels[5][0] = true

	video.Dirty = true
}
