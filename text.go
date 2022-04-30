package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	startGameText = "PRESS SPACE TO BEGIN"
	startGameLen  = len(startGameText)

	gameOverText    = "GAME OVER"
	gameOverLen     = len(gameOverText)
	gameOverOffsetY = 0

	continueText    = "PRESS SPACE TO CONTINUE"
	continueLen     = len(continueText)
	continueOffsetY = 5

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
		width:  scoreMaxWidth + padding,
		height: padding - 1,
		style:  tcell.StyleDefault,
		text:   fmt.Sprintf("Score: %d", g.score),
	})
}

func (g *Game) drawHighScore() {
	drawText(g.screen, textConfig{
		left:   scoreMaxWidth + padding - 1,
		top:    padding - 2,
		width:  highScoreMaxWidth + padding,
		height: padding - 1,
		style:  tcell.StyleDefault,
		text:   fmt.Sprintf("High Score: %d", g.highScore),
	})
}

func (g *Game) drawStartGameText() {
	drawText(
		g.screen,
		textConfig{
			left:   g.board.width/2 - continueLen/2,
			top:    g.board.height / 2,
			width:  startGameLen,
			height: g.board.height/2 + 1,
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

func (g *Game) drawPauseText() {
	drawText(
		g.screen,
		textConfig{
			left:   g.board.width/2 - pauseLen/2,
			top:    g.board.height / 2,
			width:  pauseLen,
			height: g.board.height/2 + 1,
			style:  tcell.StyleDefault,
			text:   pauseText,
		},
	)
}
