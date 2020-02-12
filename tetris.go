package main

import (
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var dropTime int = 1000 // ms
var xScale = 2
var yScale = 1
var boardSizeX = 10
var boardSizeY = 20

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func drawHints(col termbox.Attribute, left, top int) {
	x, y := left-22, top+2
	tbPrint(x, y, col, col, "  \u2190 : Move Left")
	y++
	tbPrint(x, y, col, col, "  \u2192 : Move Right")
	y++
	tbPrint(x, y, col, col, "  \u2191 : Rotate")
	y++
	tbPrint(x, y, col, col, "  \u2193 : Drop")
	y++
	tbPrint(x, y, col, col, "ESC : Quit")

	tbPrint(left+4, top+boardSizeY*yScale+3, col, col, "Play Tetris!")
}

func drawBoardFrame(col termbox.Attribute, left, top int) {
	l := left
	r := left + boardSizeX*xScale - 1
	t := top
	b := top + boardSizeY*yScale - 1

	for y := t; y <= b; y++ {
		tbPrint(l-2, y, col, col, "<|")
		tbPrint(r+1, y, col, col, "|>")

		for x := l; x <= r; x += 2 {
			tbPrint(x, y, col, col, ".")
		}
	}

	tbPrint(l-1, b+1, col, col, "+")
	tbPrint(r+1, b+1, col, col, "+")

	for x := l; x <= r; x += 2 {
		tbPrint(x, b+1, col, col, "--")
		tbPrint(x, b+2, col, col, "\\/")
	}
}

func draw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	w, h := termbox.Size()

	centerY := h / 2
	centerX := w / 2

	boardLeft := centerX - boardSizeX*xScale/2
	boardTop := centerY - boardSizeY*yScale/2 - 2

	drawBoardFrame(coldef, boardLeft, boardTop)
	drawHints(coldef, boardLeft, boardTop)

	termbox.Flush()
}

func onKeyDown() {

}

func onKeyLeft() {

}

func onKeyRight() {

}

func onKeyUp() {

}

func onDrop() {

}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	dropTick := time.NewTicker(time.Duration(dropTime) * time.Millisecond)

	draw()

mainloop:
	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					break mainloop
				case termbox.KeyArrowDown:
					onKeyDown()
				case termbox.KeyArrowLeft:
					onKeyLeft()
				case termbox.KeyArrowRight:
					onKeyRight()
				case termbox.KeyArrowUp:
					onKeyUp()
				default:
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		case <-dropTick.C:
			onDrop()
		}
		draw()
	}
}
