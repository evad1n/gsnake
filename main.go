package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
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

	if speed <= 0 {
		fmt.Printf("That's a negative speed right there sir: %f\n", speed)
		os.Exit(1)
	}

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

const highScoreFile = "highscore.txt"

// Read from file
func loadHighScore() int {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return 0
	}

	snakeDir := filepath.Join(cfgDir, "gsnake")

	err = os.MkdirAll(snakeDir, 0777)
	if err != nil {
		return 0
	}

	fullPath := filepath.Join(snakeDir, highScoreFile)

	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return 0
	}

	highScore, err := strconv.Atoi(string(raw))
	if err != nil {
		return 0
	}

	return highScore
}

// Will silently fail; not handling errors anyway...
func writeHighScore(score int) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	snakeDir := filepath.Join(cfgDir, "gsnake")

	err = os.MkdirAll(snakeDir, 0777)
	if err != nil {
		return
	}

	fullPath := filepath.Join(snakeDir, highScoreFile)

	os.WriteFile(fullPath, []byte(strconv.Itoa(score)), 0666)
}
