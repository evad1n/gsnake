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
		snakes []*Snake

		events chan GameEvent
		opts   GameOpts

		speed float64

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
	b := NewBoard(screen, opts.MaxBoardSize)

	// Calc speed based on board size
	speed := baseSpeed * float64(b.width) / 20.0

	g := &Game{
		board:  b,
		screen: screen,
		opts:   opts,
		speed:  speed,
	}

	return g
}

func (g *Game) Start() {
	g.Restart()

	for {
		if !g.over && !g.paused {
			s := g.speed * g.opts.SpeedMultiplier
			// Faster based on snake length
			s *= 1.0 + (float64(g.snakes[0].length) / 100.0)
			s *= 100.0

			s = 10000.0 / s

			mils := time.Duration(time.Millisecond * time.Duration(s))
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
		g.TurnSnakes(Up)
	case ev.Key() == tcell.KeyRight || ev.Rune() == 'd':
		g.TurnSnakes(Right)
	case ev.Key() == tcell.KeyDown || ev.Rune() == 's':
		g.TurnSnakes(Down)
	case ev.Key() == tcell.KeyLeft || ev.Rune() == 'a':
		g.TurnSnakes(Left)
	}
}

func (g *Game) TurnSnakes(dir int) {
	for _, s := range g.snakes {
		s.Turn(dir)
	}
}

func (g *Game) Update() {
	for _, s := range g.snakes {
		s.Move()
	}

	// Wall collision
	if g.CollidesWall() {
		g.over = true
		return
	}
	// Self collision
	for _, s := range g.snakes {
		if s.CollideSelf() {
			g.over = true
			return
		}
	}

	// Fruit collision
	for _, s := range g.snakes {
		if g.board.fruit.Collides(s.head.Point) {
			g.score++
			s.Grow(g.board.fruit)

			g.NewFruit()
			break
		}
	}
}

func (g *Game) SnakesCollideWithPoint(p Point) bool {
	for _, s := range g.snakes {
		if s.Collides(p) {
			return true
		}
	}
	return false
}

// Don't overlap with snake current positions
func (g *Game) NewFruit() {
	g.board.NewFruit()
	for g.SnakesCollideWithPoint(g.board.fruit) {
		g.board.NewFruit()
	}
}

func (g *Game) Draw() {
	for _, s := range g.snakes {
		s.Draw(g.screen)
	}
	g.board.Draw(g.screen)

	// Score
	drawText(g.screen, 0, 0, 10, 0, tcell.StyleDefault, fmt.Sprintf("Score: %d", g.score))

	if g.over {
		g.drawGameOverText()
	}
}

func (g *Game) Restart() {
	g.snakes = []*Snake{
		NewSnake(g.board.Midpoint()),
	}
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
	for _, s := range g.snakes {
		p := s.head

		if p.x < g.board.x ||
			p.x >= g.board.x+g.board.width ||
			p.y < g.board.y ||
			p.y >= g.board.y+g.board.height {
			return true
		}
	}

	return false
}

const (
	gameOverLen     = len("GAME OVER")
	gameOverOffsetY = 0
	continueLen     = len("PRESS SPACE TO CONTINUE")
	continueOffsetY = 5
)

func (g *Game) drawGameOverText() {
	drawText(
		g.screen,
		g.board.width/2-gameOverLen/2,
		g.board.height/2+gameOverOffsetY,
		gameOverLen,
		g.board.height/2+gameOverOffsetY+1,
		tcell.StyleDefault,
		"GAME OVER",
	)

	drawText(
		g.screen,
		g.board.width/2-continueLen/2,
		g.board.height/2+continueOffsetY,
		continueLen,
		g.board.height/2+continueOffsetY+1,
		tcell.StyleDefault,
		"PRESS SPACE TO CONTINUE",
	)
}
