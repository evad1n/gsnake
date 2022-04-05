package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type (
	Snake struct {
		length int
		head   *Cell
	}

	Cell struct {
		next *Cell
		Point
	}
)

const (
	Up int = iota
	Right
	Down
	Left

	startLength = 4
)

func NewSnake(b Board) *Snake {
	s := &Snake{
		length: startLength,
	}

	s.head = &Cell{
		Point: b.Midpoint(),
	}
	prev := s.head

	for i := s.length; i > 0; i-- {
		c := &Cell{
			Point: Point{
				x: prev.x,
				y: prev.y - 1,
			},
		}
		prev.next = c
		fmt.Println(c)
	}

	return s
}

// TODO: diff chars for directions
func (s Snake) Draw(screen tcell.Screen) {
	for c := s.head; c != nil; c = c.next {
		screen.SetContent(c.x, c.y, 'X', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorGreen))
	}
}

func (s *Snake) Move(direction int) {
	s.head.Move()
	switch direction {
	case Up:
		s.head.y++
	case Right:
		s.head.x++
	case Down:
		s.head.y--
	case Left:
		s.head.x--
	}
}

// Replace positions
func (c *Cell) Move() {
	if c.next != nil {
		c.next.Point = c.Point
		c.next.Move()
	}
}
