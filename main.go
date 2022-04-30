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
	speed       float64
	maxSize     int
	wrap        bool
	doubleSnake bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Float64Var(&speed, "speed", 1.0, "Base speed multiplier")
	flag.Float64Var(&speed, "s", 1.0, "Base speed multiplier")
	flag.IntVar(&maxSize, "size", 40, "Optional max board size")
	flag.BoolVar(&wrap, "wrap", false, "Wrap around screen")
	flag.BoolVar(&wrap, "w", false, "Wrap around screen")
	flag.BoolVar(&doubleSnake, "d", false, "Play with 2 independent snakes on the same board")
	flag.BoolVar(&doubleSnake, "double", false, "Play with 2 independent snakes on the same board")

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

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	game := NewGame(s, GameOpts{
		SpeedMultiplier: speed,
		MaxBoardSize:    maxSize,
		Wrap:            wrap,
		DoubleSnake:     doubleSnake,
	})

	go game.Start()

	// Event loop
	for {
		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				quit()
			}
			game.Event(ev)
		}
	}
}
