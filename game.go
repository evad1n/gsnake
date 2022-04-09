package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

type (
	Game struct {
		screen tcell.Screen
		board  *Board
		snake  *Snake

		events chan GameEvent
		opts   GameOpts

		score  int
		over   bool
		paused bool
	}

	GameEvent struct {
		Type string
	}

	GameOpts struct {
		SpeedMultiplier float64
		MaxBoardSize    int
	}
)

const baseSpeed float64 = 1.0

func NewGame(screen tcell.Screen, opts GameOpts) *Game {
	// Create board and snake
	b := NewBoard(screen, tcell.ColorRed)

	g := &Game{
		board:  b,
		screen: screen,
		opts:   opts,
	}

	return g
}

func (g *Game) Start() {
	g.Restart()

	for {
		if !g.over && !g.paused {
			x := baseSpeed / g.opts.SpeedMultiplier
			mils := time.Duration(time.Millisecond * time.Duration(100.0*x))
			time.Sleep(mils)
			g.Tick()
		}
	}
}

func (g *Game) Tick() {
	g.screen.Clear()

	g.Update()
	g.Draw()

	g.screen.Show()
}

func (g *Game) Event(ev *tcell.EventKey) {
	switch {
	case ev.Rune() == 'r':
		g.Restart()
	case ev.Rune() == ' ':
		if g.over {
			g.Restart()
		}
	case ev.Rune() == 'p':
		g.TogglePause()
	case ev.Key() == tcell.KeyUp || ev.Rune() == 'w':
		g.snake.Turn(Up)
	case ev.Key() == tcell.KeyRight || ev.Rune() == 'd':
		g.snake.Turn(Right)
	case ev.Key() == tcell.KeyDown || ev.Rune() == 's':
		g.snake.Turn(Down)
	case ev.Key() == tcell.KeyLeft || ev.Rune() == 'a':
		g.snake.Turn(Left)
	}
}

func (g *Game) Update() {
	g.snake.Move()

	// Wall collision
	// Self collision
	if g.CollidesWall() || g.snake.CollideSelf() {
		g.over = true

	}

	// Fruit collision
	if g.board.fruit.Collides(g.snake.head.Point) {
		g.score++
		g.snake.Grow(g.board.fruit)

		g.NewFruit()
	}
}

// Don't overlap with snake current positions
func (g *Game) NewFruit() {
	g.board.NewFruit()
	for g.snake.Collides(g.board.fruit) {
		g.board.NewFruit()
	}
}

func (g *Game) Draw() {
	g.snake.Draw(g.screen)
	g.board.Draw(g.screen)

	// Score
	drawText(g.screen, 0, 0, 10, 0, tcell.StyleDefault, fmt.Sprintf("Score: %d", g.score))
}

func (g *Game) Restart() {
	g.snake = NewSnake(*g.board)
	g.score = 0

	g.NewFruit()
	g.over = false
}

func (g *Game) Over() bool {
	return g.over
}

func (g *Game) TogglePause() {
	g.paused = !g.paused
}

func (g *Game) CollidesWall() bool {
	p := g.snake.head

	return p.x < g.board.x ||
		p.x >= g.board.x+g.board.width ||
		p.y < g.board.y ||
		p.y >= g.board.y+g.board.height
}
