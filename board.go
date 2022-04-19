package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type (
	Board struct {
		width  int
		height int
		Point

		style tcell.Style

		fruit Point
	}

	Point struct {
		x int
		y int
	}
)

const (
	fruitRune rune = 'X'
)

func NewBoard(screen tcell.Screen, padding int, maxSize int) *Board {
	w, h := screen.Size()

	if h > maxSize {
		h = maxSize
	}

	// height > width
	idealWidth := h * 2
	if idealWidth > w {
		idealWidth = w
	}

	return &Board{
		width:  idealWidth - padding*2,
		height: h - padding*2,
		Point:  Point{padding, padding},
		style:  tcell.StyleDefault.Background(tcell.ColorWhite),
	}
}

func (b *Board) Draw(screen tcell.Screen) {
	// Vertical walls
	for row := b.y; row <= b.y+b.height; row++ {
		screen.SetContent(b.x-1, row, ' ', nil, b.style)
		screen.SetContent(b.x+b.width, row, ' ', nil, b.style)
	}
	// Top
	for col := b.x - 1; col < b.x+b.width+1; col++ {
		screen.SetContent(col, b.y-1, ' ', nil, b.style)
	}
	// Bottom
	for col := b.x; col < b.x+b.width; col++ {
		screen.SetContent(col, b.y+b.height, ' ', nil, b.style)
	}

	// Draw fruit
	screen.SetContent(b.fruit.x, b.fruit.y, fruitRune, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func (b *Board) NewFruit() {
	x := rand.Intn(b.width)
	y := rand.Intn(b.height)

	b.fruit = Point{b.x + x, b.y + y}
}

func (b *Board) Midpoint() Point {
	return Point{
		x: b.x + (b.width / 2),
		y: b.y + (b.height / 2),
	}
}

func (p Point) Collides(other Point) bool {
	return p.x == other.x && p.y == other.y
}

// -1 default return if points aren't one block away in cardinal directions
func (p Point) DirTo(other Point) int {
	y := p.y - other.y
	x := p.x - other.x

	switch {
	case y == -1:
		return Down
	case y == 1:
		return Up
	case x == -1:
		return Right
	case x == 1:
		return Left

	default:
		return -1
	}
}

func (p Point) Equals(other Point) bool {
	return p.x == other.x && p.y == other.y
}
