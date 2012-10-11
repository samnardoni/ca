package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"github.com/jteeuwen/glfw"
	"math/rand"
	"os"
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

const (
	Title  = "Cellular Automata"
	Width  = 480
	Height = 480
	Size   = 50
)

type GridBuf struct {
	buf     [][][]bool
	current int
}

func NewGridBuf(size int) *GridBuf {
	buf := make([][][]bool, 2)

	for b := 0; b < 2; b++ {
		buf[b] = make([][]bool, size)

		for i := 0; i < size; i++ {
			buf[b][i] = make([]bool, size)
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
	running bool    = true
	counter Counter = Counter{hit: 3, current: 0}
)

func init() {
	grid = NewGridBuf(Size)
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

	initGL()
	initGrid()

	for running && glfw.WindowParam(glfw.Opened) == 1 {
		update()
		drawScene()
	}
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = false
	}
}

func initGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

func initGrid() {
	for i := 0; i < (Size*Size)/2; i++ {
		rx := rand.Int31n(Size - 1)
		ry := rand.Int31n(Size - 1)
		grid.Front()[ry][rx] = true
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
			sum := 0
			if f[y-1][x-1] {
				sum++
			}
			if f[y-1][x] {
				sum++
			}
			if f[y-1][x+1] {
				sum++
			}
			if f[y][x-1] {
				sum++
			}
			if f[y][x+1] {
				sum++
			}
			if f[y+1][x-1] {
				sum++
			}
			if f[y+1][x] {
				sum++
			}
			if f[y+1][x+1] {
				sum++
			}

			if sum == 2 {
				b[y][x] = f[y][x]
			} else if sum == 3 {
				b[y][x] = true
			} else {
				b[y][x] = false
			}
		}
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for y := 1; y < Size-1; y++ {
		for x := 1; x < Size-1; x++ {
			gl.LoadIdentity()
			gl.Scalef(0.85/float32(Size), 0.85/float32(Size), 1)
			gl.Translatef(float32(x)-Size/2, float32(y)-Size/2, -1)

			if grid.Front()[y][x] {
				gl.Color3f(0.0, 0.0, 0.0)
			} else {
				gl.Color3f(1.0, 1.0, 1.0)
			}

			gl.Begin(gl.QUADS)
			gl.Vertex3f(0, 1, 0)
			gl.Vertex3f(1, 1, 0)
			gl.Vertex3f(1, 0, 0)
			gl.Vertex3f(0, 0, 0)
			gl.End()
		}
	}

	glfw.SwapBuffers()
}
