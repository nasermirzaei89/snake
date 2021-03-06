package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math/rand"
	"time"
)

const (
	screenWidth  = 48
	screenHeight = 27
	tileWidth    = 8
)

const (
	directionRight = iota
	directionUp
	directionLeft
	directionDown
)

const waitTicks = 10

type point struct {
	x, y int
}

type game struct {
	rnd              *rand.Rand
	snake            []point
	direction        int
	directionHandled bool
	seedSowed        bool
	gameOver         bool
	gameStarted      bool
	tick             int
	seed             point
	walls            []point
	highScore        int
}

func (g *game) seeding() {
	if len(g.snake)+len(g.walls) == screenWidth*screenHeight {
		g.gameOver = true

		return
	}

L1:
	for {
		g.seed.x = g.rnd.Intn(screenWidth)
		g.seed.y = g.rnd.Intn(screenHeight)

		for i := range g.snake {
			if g.seed.x == g.snake[i].x && g.seed.y == g.snake[i].y {
				continue L1
			}
		}

		for i := range g.walls {
			if g.seed.x == g.walls[i].x && g.seed.y == g.walls[i].y {
				continue L1
			}
		}

		break
	}

	g.seedSowed = true
}

func (g *game) score() int {
	if g.snake == nil {
		return 0
	}

	return len(g.snake) - 2
}

func (g *game) Update() error {
	if g.rnd == nil {
		g.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// init walls
	g.walls = make([]point, 0)
	for i := 0; i < screenWidth; i++ {
		g.walls = append(
			g.walls,
			point{
				x: i,
				y: 0,
			},
			point{
				x: i,
				y: screenHeight - 1,
			},
		)
	}

	// init snake
	if g.snake == nil {
		g.snake = []point{
			{x: screenWidth / 2, y: screenHeight / 2},
		}

		g.snake = append(g.snake, point{x: g.snake[0].x - 1, y: g.snake[0].y})
	}

	if !g.seedSowed && !g.gameOver {
		g.seeding()
	}

	if g.gameOver && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.snake = nil
		g.seedSowed = false
		g.gameOver = false
		g.gameStarted = false
		g.directionHandled = false
	}

	// handle input
	if !g.directionHandled && inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.direction++
		g.directionHandled = true
		g.gameStarted = true
	}

	if !g.directionHandled && inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.direction--
		g.directionHandled = true
		g.gameStarted = true
	}

	if !g.directionHandled && g.direction != directionLeft && inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.direction = directionRight
		g.directionHandled = true
		g.gameStarted = true
	}

	if !g.directionHandled && g.direction != directionDown && inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.direction = directionUp
		g.directionHandled = true
		g.gameStarted = true
	}

	if !g.directionHandled && g.direction != directionRight && inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.direction = directionLeft
		g.directionHandled = true
		g.gameStarted = true
	}

	if !g.directionHandled && g.direction != directionUp && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.direction = directionDown
		g.directionHandled = true
		g.gameStarted = true
	}

	g.direction = (g.direction + 4) % 4

	// step
	if g.gameStarted && !g.gameOver {
		if g.tick%waitTicks == 0 {
			g.tick = 0

			head := g.snake[0]

			switch g.direction {
			case directionLeft:
				head.x--
			case directionUp:
				head.y--
			case directionRight:
				head.x++
			case directionDown:
				head.y++
			}

			// screen wrap
			head.x = (head.x + screenWidth) % screenWidth
			head.y = (head.y + screenHeight) % screenHeight

			g.snake = append([]point{head}, g.snake...)

			// check eat
			if head.x == g.seed.x && head.y == g.seed.y {
				g.seedSowed = false

				if s := g.score(); s > g.highScore {
					g.highScore = s
				}
			} else {
				g.snake = g.snake[:len(g.snake)-1]
			}

			// check bite self
			for i := 1; i < len(g.snake); i++ {
				if head.x == g.snake[i].x && head.y == g.snake[i].y {
					g.gameOver = true
				}
			}

			// check crash
			for i := 0; i < len(g.walls); i++ {
				if head.x == g.walls[i].x && head.y == g.walls[i].y {
					g.gameOver = true
				}
			}

			g.directionHandled = false
		}

		g.tick++
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	for i := range g.walls {
		ebitenutil.DrawRect(screen, float64(g.walls[i].x*tileWidth), float64(g.walls[i].y*tileWidth), tileWidth, tileWidth, color.RGBA{
			R: 127,
			G: 127,
			B: 127,
			A: 255,
		})
	}

	if g.seedSowed {
		ebitenutil.DrawRect(screen, float64(g.seed.x*tileWidth), float64(g.seed.y*tileWidth), tileWidth, tileWidth, color.RGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 255,
		})
	}

	for i := len(g.snake) - 1; i >= 0; i-- {
		c := color.RGBA{
			R: 0,
			G: 255,
			B: 0,
			A: 255,
		}

		if g.gameOver && i == 0 {
			// read head
			c = color.RGBA{
				R: 255,
				G: 0,
				B: 0,
				A: 255,
			}
		}

		ebitenutil.DrawRect(screen, float64(g.snake[i].x*tileWidth), float64(g.snake[i].y*tileWidth), tileWidth, tileWidth, c)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(" [ Score: %d ] [ High Score: %d ]", g.score(), g.highScore))
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth * tileWidth, screenHeight * tileWidth
}

var _ ebiten.Game = new(game)
