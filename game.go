package main

import (
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

		score     int
		highScore int

		over    bool
		paused  bool
		started bool
	}

	GameEvent struct {
		Type string
	}

	GameOpts struct {
		SpeedMultiplier float64
		MaxBoardSize    int
		Wrap            bool
		DoubleSnake     bool
	}
)

const (
	baseSpeed float64 = 1.0
	padding   int     = 3
)

func NewGame(screen tcell.Screen, opts GameOpts) *Game {
	// Create board and snake
	b := NewBoard(screen, padding, opts.MaxBoardSize)

	// Calc speed based on board size
	speed := baseSpeed * float64(b.width) / 20.0

	g := &Game{
		board:     b,
		screen:    screen,
		opts:      opts,
		speed:     speed,
		highScore: loadHighScore(),
	}

	return g
}

func (g *Game) Start() {
	g.screen.Clear()
	g.Draw()
	g.screen.Show()

	for {
		if g.IsPlayingState() {
			s := g.speed * g.opts.SpeedMultiplier
			// Faster based on snake length
			s *= 0.6 + (float64(g.snakes[0].length) / 100.0)
			// This calculation is psychotic
			s *= 100.0

			s = 10000.0 / s

			mils := time.Duration(time.Millisecond * time.Duration(s))
			time.Sleep(mils)
			g.Tick()
		}
	}
}

func (g *Game) IsPlayingState() bool {
	return g.started && !g.over && !g.paused
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
	// Space
	case ev.Rune() == ' ':
		if g.over || !g.started {
			g.Restart()
		} else {
			g.TogglePause()
		}
	case ev.Rune() == 'n' && g.paused:
		g.Tick()
	default:
		// Let turns happen while paused for frame by frame
		if g.started && !g.over {
			g.TurnSnakes(ev)
		}
	}
}

func (g *Game) Update() {
	for _, s := range g.snakes {
		s.Move()
		s.UpdateDraw()
	}

	// Wall collision
	if g.opts.Wrap {
		g.WrapSnakes()
	} else {
		if g.CollidesWall() {
			g.over = true
			return
		}
	}

	// Self collision
	for _, s := range g.snakes {
		if s.CollideSelf() {
			g.over = true
			return
		}
	}

	// Multi-snake collision
	if g.opts.DoubleSnake {
		for i := range g.snakes {
			for j := range g.snakes {
				if i == j {
					continue
				}
				if g.snakes[i].Collides(g.snakes[j].head.Point) {
					g.over = true
					return
				}

			}
		}
	}

	// Fruit collision
	for _, s := range g.snakes {
		if g.board.fruit.Collides(s.head.Point) {
			g.score++
			s.Grow(g.board.fruit)

			if g.score > g.highScore {
				g.highScore = g.score
				writeHighScore(g.score)
			}

			g.NewFruit()
			break
		}
	}
}

func (g *Game) Draw() {
	for _, s := range g.snakes {
		s.Draw(g.screen)
	}

	g.board.Draw(g.screen, g.started)

	if !g.started {
		g.drawStartGameText()
	}

	if g.paused {
		g.drawPauseText()
	}

	g.drawScore()
	if g.highScore != -1 {
		g.drawHighScore()
	}

	if g.over {
		g.drawGameOverText()
	}

}

func (g *Game) TurnSnakes(ev *tcell.EventKey) {
	switch {
	case ev.Rune() == 'w':
		if g.opts.DoubleSnake {
			g.snakes[0].Turn(Up)
		} else {
			for _, s := range g.snakes {
				s.Turn(Up)
			}
		}
	case ev.Rune() == 'd':
		if g.opts.DoubleSnake {
			g.snakes[0].Turn(Right)
		} else {
			for _, s := range g.snakes {
				s.Turn(Right)
			}
		}
	case ev.Rune() == 's':
		if g.opts.DoubleSnake {
			g.snakes[0].Turn(Down)
		} else {
			for _, s := range g.snakes {
				s.Turn(Down)
			}
		}
	case ev.Rune() == 'a':
		if g.opts.DoubleSnake {
			g.snakes[0].Turn(Left)
		} else {
			for _, s := range g.snakes {
				s.Turn(Left)
			}
		}
	case ev.Key() == tcell.KeyUp:
		if g.opts.DoubleSnake {
			g.snakes[1].Turn(Up)
		} else {
			for _, s := range g.snakes {
				s.Turn(Up)
			}
		}
	case ev.Key() == tcell.KeyRight:
		if g.opts.DoubleSnake {
			g.snakes[1].Turn(Right)
		} else {
			for _, s := range g.snakes {
				s.Turn(Right)
			}
		}
	case ev.Key() == tcell.KeyDown:
		if g.opts.DoubleSnake {
			g.snakes[1].Turn(Down)
		} else {
			for _, s := range g.snakes {
				s.Turn(Down)
			}
		}
	case ev.Key() == tcell.KeyLeft:
		if g.opts.DoubleSnake {
			g.snakes[1].Turn(Left)
		} else {
			for _, s := range g.snakes {
				s.Turn(Left)
			}
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

func (g *Game) Restart() {
	if g.opts.DoubleSnake {
		mid := g.board.Midpoint()
		g.snakes = []*Snake{
			NewSnake(Point{
				x: mid.x - 2,
				y: mid.y,
			}),
			NewSnake(Point{
				x: mid.x + 2,
				y: mid.y,
			}),
		}
	} else {
		g.snakes = []*Snake{
			NewSnake(g.board.Midpoint()),
		}
	}

	g.score = 0

	g.NewFruit()
	g.started = true
	g.over = false

	g.Tick()
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

func (g *Game) WrapSnakes() {
	for _, s := range g.snakes {
		p := s.head

		switch {
		case p.x < g.board.x:
			s.head.x = g.board.x + g.board.width - 1
		case p.x >= g.board.x+g.board.width:
			s.head.x = g.board.x
		case p.y < g.board.y:
			s.head.y = g.board.y + g.board.height - 1
		case p.y >= g.board.y+g.board.height:
			s.head.y = g.board.y
		}
	}
}
