package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func draw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	w, h := termbox.Size()

	midy := h / 2
	midx := w / 2

	tbprint(midx-9, midy-1, coldef, coldef, "Let's code Tetris!")
	tbprint(midx-9, midy+1, coldef, coldef, "Press ESC to quit")

	termbox.Flush()
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	draw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			default:
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw()
	}
}
