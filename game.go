package main

import "github.com/gdamore/tcell/v2"

type (
	Game struct {
		board *Board
		snake *Snake
	}
)

func NewGame(screen tcell.Screen) *Game {
	// Create board and snake
	b := NewBoard(screen, tcell.ColorRed)
	snake := NewSnake(*b)
	return &Game{
		board: b,
		snake: snake,
	}
}

func (g *Game) Update() {
	g.snake.Move()
}
