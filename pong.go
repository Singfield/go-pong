package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  int32 = 800
	winHeight int32 = 600
)

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius int
	xv     float32
	yv     float32

	color color
}

type paddle struct {
	pos
	w     int
	h     int
	color color
}

// Go doesn't have ineritence but have composition
// instead of typing :

// type paddle struct {
// 	pos   pos
// 	w     int
// 	h     int
// 	color color
// }
// which gives : paddle.pos.x
// we tape this :
// type paddle struct {
// 	pos
// 	w     int
// 	h     int
// 	color color
// }
//which gives: paddle.x
// with that you can use receive func (methode) from pos in paddle

func (paddle *paddle) draw(pixels []byte) {

	// We set from the center the starting point to drow at the left hand corner

	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (ball *ball) draw(pixels []byte) {
	// parametric equation of a cicle : x^2 + y^2 =r^2
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func ( ball *ball) update(leftPaddle *paddle, rightPaddle *paddle){
	ball.x += ball.xv
	ball.y += ball.yv

	// handle collision

	if ball.y -float32(ball.radius) < 0 || ball.y+float32(ball.radius) > float32(winHeight){
		ball.yv = -ball.yv
	}

	if ball.x < 0 || ball.x > float32(winWidth) {
		ball.x  = float32(winWidth)/2
		ball.y = float32(winHeight)/2
	}

	if ball.x < leftPaddle.x + float32(rightPaddle.w/2){
		if ball.y > leftPaddle.y-float32(leftPaddle.h/2) && ball.y < leftPaddle.y+float32(leftPaddle.h/2) {
			ball.xv = -ball.xv
		}
	}

	if ball.x > rightPaddle.x - float32(rightPaddle.w/2) {
		if ball.y > rightPaddle.y-float32(rightPaddle.h/2) && ball.y < rightPaddle.y+float32(rightPaddle.h/2) {
			ball.xv = -ball.xv
		}
	}
}

func (paddle *paddle) update(keyState []uint8){
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -=5 
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y+=5
	}
	// need user inputs ;)
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*int(winWidth) + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}

}

func clear( pixels []byte) {
	for i := range pixels {
		pixels[i] =0
	}
}

func (paddle *paddle) aiUpdate(ball *ball){
	paddle.y = ball.y

}

func main() {

	// Added after EPO6 to address macosx issues
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Pong Sdl2 Golang", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Errorf("impossible to open a windows %v", err)
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, winWidth, winHeight)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	

	player1 := paddle{pos{50, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth)-50, 100}, 20, 100, color{255, 255, 255}}

	ball := ball{pos{300, 300}, 20, 5, 5, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	//Changed after EP06 to address MacOSX
	//OSX Requires that to continu events for windows to open and work properly

	// It's also be our game loop
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		player1.update(keyState)

		ball.update(&player1, &player2)
		player2.aiUpdate(&ball)

		player1.draw(pixels)
		ball.draw(pixels)
		player2.draw(pixels)
		
		tex.Update(nil, pixels, int(winWidth)*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}
