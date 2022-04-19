package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	speed   float64
	maxSize int
	wrap    bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Float64Var(&speed, "speed", 1.0, "Base speed multiplier")
	flag.Float64Var(&speed, "s", 1.0, "Base speed multiplier")
	flag.IntVar(&maxSize, "size", 40, "Optional max board size")
	flag.BoolVar(&wrap, "wrap", false, "Wrap around screen")
	flag.BoolVar(&wrap, "w", false, "Wrap around screen")

	flag.Parse()

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	// Event loop
	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	game := NewGame(s, GameOpts{
		SpeedMultiplier: speed,
		MaxBoardSize:    maxSize,
		Wrap:            wrap,
	})

	go game.Start()

	// Event listener
	for {
		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
			game.Event(ev)
		}
	}
}

func drawText(s tcell.Screen, left, top, width, height int, style tcell.Style, text string) {
	row := top
	col := left
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= left+width {
			row++
			col = left
		}
		if row > top+height {
			break
		}
	}
}
