package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

var (
	startGameText = "PRESS SPACE TO BEGIN"
	startGameLen  = len(startGameText)

	gameOverText = "GAME OVER"
	gameOverLen  = len(gameOverText)

	continueText = "PRESS SPACE TO CONTINUE"
	continueLen  = len(continueText)

	pauseText = "PAUSED"
	pauseLen  = len(pauseText)

	scoreMaxWidth     = 10
	highScoreMaxWidth = scoreMaxWidth + 5
)

type textConfig struct {
	left   int
	top    int
	width  int
	height int
	style  tcell.Style
	text   string
}

func drawText(s tcell.Screen, cfg textConfig) {
	row := cfg.top
	col := cfg.left

	maxCol := cfg.left + cfg.width
	maxRow := cfg.top + cfg.height

	for _, r := range []rune(cfg.text) {
		s.SetContent(col, row, r, nil, cfg.style)
		col++
		if col >= maxCol {
			row++
			col = cfg.left
		}
		if row > maxRow {
			break
		}
	}
}

// Accounts for board padding
func drawBoardText(s tcell.Screen, b *Board, cfg textConfig) {
	row := cfg.top + b.padding
	col := cfg.left + b.padding

	maxCol := cfg.left + cfg.width + b.padding
	maxRow := cfg.top + cfg.height + b.padding

	for _, r := range []rune(cfg.text) {
		s.SetContent(col, row, r, nil, cfg.style)
		col++
		if col >= maxCol {
			row++
			col = cfg.left
		}
		if row > maxRow {
			break
		}
	}
}

func (g *Game) drawScore() {
	drawText(g.screen, textConfig{
		left:   padding - 1,
		top:    padding - 2,
		width:  scoreMaxWidth,
		height: padding - 1,
		style:  tcell.StyleDefault,
		text:   fmt.Sprintf("Score: %d", g.score),
	})
}

func (g *Game) drawHighScore() {
	drawText(g.screen, textConfig{
		left:   g.board.width + g.board.padding - highScoreMaxWidth,
		top:    padding - 2,
		width:  highScoreMaxWidth,
		height: padding - 1,
		style:  tcell.StyleDefault,
		text:   fmt.Sprintf("High Score: %d", g.highScore),
	})
}

func (g *Game) drawStartGameText() {
	drawBoardText(
		g.screen,
		g.board,
		textConfig{
			left:   g.board.width/2 - startGameLen/2,
			top:    int(float32(g.board.height/2) * 0.9),
			width:  startGameLen,
			height: 2,
			style:  tcell.StyleDefault,
			text:   startGameText,
		},
	)
}

func (g *Game) drawGameOverText() {
	drawBoardText(
		g.screen,
		g.board,
		textConfig{
			left:   g.board.width/2 - gameOverLen/2,
			top:    int(float32(g.board.height/2) * 0.8),
			width:  gameOverLen,
			height: 2,
			style:  tcell.StyleDefault,
			text:   gameOverText,
		},
	)

	drawBoardText(
		g.screen,
		g.board,
		textConfig{
			left:   g.board.width/2 - continueLen/2,
			top:    int(float32(g.board.height/2)*0.8) + 2,
			width:  continueLen,
			height: 2,
			style:  tcell.StyleDefault,
			text:   continueText,
		},
	)
}

func (g *Game) drawPauseText() {
	drawBoardText(
		g.screen,
		g.board,
		textConfig{
			left:   g.board.width/2 - pauseLen/2,
			top:    int(float32(g.board.height/2) * 0.9),
			width:  pauseLen,
			height: 1,
			style:  tcell.StyleDefault,
			text:   pauseText,
		},
	)
}
