package main

import (
	"github.com/gdamore/tcell/v2"
)

type (
	Snake struct {
		length int
		head   *Cell
		dir    int

		// Can potentially be digesting multiple things if we are real long
		growPositions []Point
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

	startLength = 6
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
		prev = c
	}

	return s
}

// Super cool animation right
func (s *Snake) Grow(pos Point) {
	// s.length++
	// var c *Cell
	// // Skip to tail
	// for c = s.head; c != nil; c = c.next {
	// }
	// c.next = &Cell{
	// 	Point: Point{
	// 		x: c.x,
	// 		y: c.y,
	// 	},
	// }
}

func (s Snake) Draw(screen tcell.Screen) {
	var prev *Cell
	for c := s.head; c != nil; c = c.next {
		c.Draw(screen, prev)
		prev = c
	}
}

// Replace positions
func (c *Cell) Draw(screen tcell.Screen, prev *Cell) {
	char := 'o'
	switch {
	// Head
	case prev == nil:
		switch {
		case c.next.x < c.x:
			char = '>'
		case c.next.x > c.x:
			char = '<'
		case c.next.y < c.y:
			char = 'v'
		case c.next.y > c.y:
			char = '^'
		}
	// Tail
	case c.next == nil:
		switch {
		case prev.x < c.x:
			char = '>'
		case prev.x > c.x:
			char = '<'
		case prev.y < c.y:
			char = 'v'
		case prev.y > c.y:
			char = '^'
		}
	}

	// ╔
	// ╚
	// ╗
	// ╝

	screen.SetContent(c.x, c.y, char, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
}

func (s *Snake) Turn(direction int) {
	// Only turn 90 degress max
	diff := s.dir - direction
	if diff < 0 {
		diff = -diff
	}
	// 180
	if diff == 2 {
		return
	}
	s.dir = direction
}

func (s *Snake) Move() {
	s.head.Move()
	switch s.dir {
	case Up:
		s.head.y--
	case Right:
		s.head.x++
	case Down:
		s.head.y++
	case Left:
		s.head.x--
	}
}

// Replace positions and move up
func (c *Cell) Move() {
	if c.next != nil {
		c.next.Move()
		c.next.Point = c.Point
	}
}

// Will check for every possible cell collision
func (s *Snake) Collides(p Point) bool {
	for c := s.head; c != nil; c = c.next {
		if c.Point.Collides(p) {
			return true
		}
	}
	return false
}
