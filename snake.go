package main

import (
	"github.com/gdamore/tcell/v2"
)

type (
	Snake struct {
		length int
		head   *Cell
		dir    int

		// Otherwise you can turn 180 with multiple inputs per tick
		pendingTurn int

		// Decide which runes to draw, then transform
		pendingDraw []rune

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

func NewSnake(start Point) *Snake {
	s := &Snake{
		length: startLength,
	}

	s.head = &Cell{
		Point: start,
	}
	prev := s.head

	for i := s.length - 1; i > 0; i-- {
		c := &Cell{
			Point: Point{
				x: prev.x,
				y: prev.y + 1,
			},
		}
		prev.next = c
		prev = c
	}

	return s
}

// Pending growth to make it look cooler
func (s *Snake) Grow(pos Point) {
	s.growPositions = append(s.growPositions, pos)
}

func (s Snake) Draw(screen tcell.Screen) {
	i := 0
	for c := s.head; c != nil; c = c.next {
		screen.SetContent(c.x, c.y, s.pendingDraw[i], nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		i++
	}
	s.pendingDraw = nil
}

func (s *Snake) UpdateDraw() {
	s.pendingDraw = make([]rune, s.length)
	var prev *Cell
	i := 0
	for c := s.head; c != nil; c = c.next {
		c.PendingDraw(s.pendingDraw, i, prev, s.CellIsGrowthPosition(c))
		prev = c
		i++
	}
}

// Replace positions
func (c *Cell) PendingDraw(pendingDraw []rune, index int, prev *Cell, isGrow bool) {
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
			char = '-'
		case prev.x > c.x:
			char = '-'
		case prev.y < c.y:
			char = '|'
		case prev.y > c.y:
			char = '|'
		}
	case isGrow:
		char = 'O'
	default:
		dirToPrev := c.Point.DirTo(prev.Point)
		dirToNext := c.Point.DirTo(c.next.Point)

		switch {
		case (dirToPrev == Up || dirToPrev == Down) &&
			(dirToNext == Up || dirToNext == Down):
			char = '|'
		case (dirToPrev == Left || dirToPrev == Right) &&
			(dirToNext == Left || dirToNext == Right):
			char = '-'
		case dirToPrev == Left && dirToNext == Up ||
			dirToPrev == Up && dirToNext == Left:
			char = '┘'
		case dirToPrev == Right && dirToNext == Up ||
			dirToPrev == Up && dirToNext == Right:
			char = '└'
		case dirToPrev == Left && dirToNext == Down ||
			dirToPrev == Down && dirToNext == Left:
			char = '┐'
		case dirToPrev == Right && dirToNext == Down ||
			dirToPrev == Down && dirToNext == Right:
			char = '┌'
		}
	}

	pendingDraw[index] = char
}

func (s *Snake) CellIsGrowthPosition(c *Cell) bool {
	for _, p := range s.growPositions {
		if p.Equals(c.Point) {
			return true
		}
	}

	return false
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
	s.pendingTurn = direction
}

func (s *Snake) Move() {
	// Add pending growth
	s.UpdateGrowth()

	s.dir = s.pendingTurn
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

func (s *Snake) Tail() *Cell {
	c := s.head
	for {
		if c.next == nil {
			return c
		}
		c = c.next
	}
}

func (s *Snake) UpdateGrowth() {
	tail := s.Tail()

	for i, p := range s.growPositions {
		if tail.Equals(p) {
			s.length++
			tail.next = &Cell{
				Point: p,
			}

			// Remove from grow positions
			s.growPositions = append(s.growPositions[:i], s.growPositions[i+1:]...)
		}
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

// Just check if the head is colliding
func (s *Snake) CollideSelf() bool {
	for c := s.head.next; c != nil; c = c.next {
		if c.Point.Collides(s.head.Point) {
			return true
		}
	}

	return false
}
