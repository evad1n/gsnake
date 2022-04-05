package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type (
	Game struct {
		board  *Board
		snake  *Snake
		screen tcell.Screen
	}
)

const baseSpeed time.Duration = time.Millisecond * 200

func NewGame(screen tcell.Screen) *Game {
	// Create board and snake
	b := NewBoard(screen, tcell.ColorRed)
	snake := NewSnake(*b)

	return &Game{
		board:  b,
		snake:  snake,
		screen: screen,
	}
}

func (g *Game) Start() {
	g.NewFruit()

	for !g.Over() {
		g.screen.Clear()

		g.Update()
		g.Draw()

		g.screen.Show()
		time.Sleep(baseSpeed)
	}

	// Some stuff
}

func (g *Game) Event(ev *tcell.EventKey) {
	switch {
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

	// Fruit collision
	if g.board.fruit.Collides(g.snake.head.Point) {
		g.snake.Grow(g.board.fruit)

		g.NewFruit()
	}
}

// Don't overlap with snake current positions
func (g *Game) NewFruit() {
	for g.snake.Collides(g.board.fruit) {
		g.board.NewFruit()
	}
}

func (g *Game) Draw() {
	g.board.Draw(g.screen)
	g.snake.Draw(g.screen)
}

func (g *Game) Over() bool {
	return false
}
