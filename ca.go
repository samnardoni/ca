package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"math/rand"
	"os"
)

const (
	Title  = "Cellular Automata"
	Size   = 200
	Width  = Size * 3
	Height = Size * 3
)

type Counter struct {
	hit, current int
}

func (c *Counter) Tick() bool {
	c.current++
	if c.current%c.hit == 0 {
		return true
	}
	return false
}

type GridBuf struct {
	buf     [2][][]bool
	current int
}

func NewGridBuf(size int) *GridBuf {
	var buf [2][][]bool

	for b := 0; b < 2; b++ {
		buf[b] = make([][]bool, size)

		for i := 0; i < size; i++ {
			buf[b][i] = make([]bool, size)
		}
	}

	return &GridBuf{buf: buf, current: 0}
}

func (g *GridBuf) Front() [][]bool {
	return g.buf[g.current]
}

func (g *GridBuf) Back() [][]bool {
	return g.buf[1-g.current]
}

func (g *GridBuf) Swap() {
	g.current = 1 - g.current
}

var (
	grid    *GridBuf
	pixels          = make([]uint8, Size*Size*3)
	running bool    = true
	counter Counter = Counter{hit: 3, current: 0}
)

func init() {
	grid = NewGridBuf(Size)

	// Randomise grid
	for i := 0; i < (Size*Size)/4; i++ {
		rx := rand.Int31n(Size-2) + 1
		ry := rand.Int31n(Size-2) + 1
		grid.Front()[ry][rx] = true
	}
}

func main() {
	var err error

	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.Terminate()

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.TEXTURE_2D)

	for running && glfw.WindowParam(glfw.Opened) == 1 {
		update()
		draw()
	}
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-1, 1, -1, 1, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Disable(gl.DEPTH_TEST)
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = false
	}
}

func update() {
	if hit := counter.Tick(); !hit {
		return
	}

	f := grid.Front()
	b := grid.Back()

	defer grid.Swap()

	for y := 1; y < Size-1; y++ {
		for x := 1; x < Size-1; x++ {

			neighbours := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if !(dx == 0 && dy == 0) && f[y+dy][x+dx] {
						neighbours++
					}
				}
			}

			if neighbours == 2 {
				b[y][x] = f[y][x]
			} else if neighbours == 3 {
				b[y][x] = true
			} else {
				b[y][x] = false
			}

		}
	}
}

func draw() {
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {

			var color uint8
			if grid.Front()[y][x] {
				color = 0x00
			} else {
				color = 0xFF
			}

			pixels[(y*Size+x)*3+0] = color
			pixels[(y*Size+x)*3+1] = color
			pixels[(y*Size+x)*3+2] = color

		}
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, 3, Size, Size, 0, gl.RGB, gl.UNSIGNED_BYTE, pixels)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, 0.0)
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, 0.0)
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(1.0, 1.0, 0.0)
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(-1.0, 1.0, 0.0)
	gl.End()

	glfw.SwapBuffers()
}
