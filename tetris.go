package main

import (
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var dropTime = 1000 // ms
var info = "Event Test"

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

	tbprint(midx-9, midy-1, coldef, coldef, info)
	tbprint(midx-9, midy+1, coldef, coldef, "Press ESC to quit")

	termbox.Flush()
}

func onKeyDown() {
	info = "Down"
}

func onKeyLeft() {
	info = "Left"
}

func onKeyRight() {
	info = "Right"
}

func onKeyUp() {
	info = "Up"
}

func onDrop() {
	info = "Drop"
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
