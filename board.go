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

const padding = 2

func NewBoard(screen tcell.Screen, color tcell.Color) *Board {
	w, h := screen.Size()

	return &Board{
		width:  w - padding*2,
		height: h - padding*2,
		Point:  Point{padding, padding},
		style:  tcell.StyleDefault.Foreground(color),
	}
}

func (b *Board) Draw(screen tcell.Screen) {
	// Draw borders
	for col := b.x; col <= b.width; col++ {
		screen.SetContent(col, b.y, tcell.RuneHLine, nil, b.style)
		screen.SetContent(col, b.y+b.height, tcell.RuneHLine, nil, b.style)
	}
	for row := b.y + 1; row <= b.height+1; row++ {
		screen.SetContent(b.x, row, tcell.RuneVLine, nil, b.style)
		screen.SetContent(b.x+b.width, row, tcell.RuneVLine, nil, b.style)
	}
}

func (b *Board) NewFruit() {
	x := rand.Intn(b.width)
	y := rand.Intn(b.height)

	b.fruit = Point{x, y}
}

func (b *Board) Midpoint() Point {
	return Point{
		x: b.x + (b.width / 2),
		y: b.y + (b.height / 2),
	}
}
