package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	startGameText    = "PRESS SPACE TO BEGIN"
	startGameLen     = len(startGameText)
	startGameOffsetY = 0

	gameOverText    = "GAME OVER"
	gameOverLen     = len(gameOverText)
	gameOverOffsetY = 0

	continueText    = "PRESS SPACE TO CONTINUE"
	continueLen     = len(continueText)
	continueOffsetY = 5
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
	for _, r := range []rune(cfg.text) {
		s.SetContent(col, row, r, nil, cfg.style)
		col++
		if col >= cfg.left+cfg.width {
			row++
			col = cfg.left
		}
		if row > cfg.top+cfg.height {
			break
		}
	}
}

func (g *Game) drawScore() {
	drawText(g.screen, textConfig{
		left:   padding - 1,
		top:    padding - 2,
		width:  10 + padding,
		height: padding - 1,
		style:  tcell.StyleDefault,
		text:   fmt.Sprintf("Score: %d", g.score),
	})
}

func (g *Game) drawStartGameText() {
	drawText(
		g.screen,
		textConfig{
			left:   g.board.width/2 - continueLen/2,
			top:    g.board.height/2 + startGameOffsetY,
			width:  startGameLen,
			height: g.board.height/2 + startGameOffsetY + 1,
			style:  tcell.StyleDefault,
			text:   startGameText,
		},
	)
}

func (g *Game) drawGameOverText() {
	drawText(
		g.screen,
		textConfig{
			left:   g.board.width/2 - gameOverLen/2,
			top:    g.board.height/2 + gameOverOffsetY,
			width:  gameOverLen,
			height: g.board.height/2 + gameOverOffsetY + 1,
			style:  tcell.StyleDefault,
			text:   gameOverText,
		},
	)

	drawText(
		g.screen,
		textConfig{
			left:   g.board.width/2 - continueLen/2,
			top:    g.board.height/2 + continueOffsetY,
			width:  continueLen,
			height: g.board.height/2 + continueOffsetY + 1,
			style:  tcell.StyleDefault,
			text:   continueText,
		},
	)
}
