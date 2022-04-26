package main

import "github.com/gdamore/tcell/v2"

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
