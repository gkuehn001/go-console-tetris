package main

import (
	"math/rand"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const initialDropInterval = 1000 // ms
const hardDropInterval = 20      // ms
const xScale = 2
const yScale = 1
const boardSizeX = 10
const boardSizeY = 20

type dropResult int

const (
	none dropResult = 0
	fix  dropResult = 1
	full dropResult = 2
)

var tetrominoX = [7][16]int{
	{0, 1, 2, 3, 2, 2, 2, 2, 0, 1, 2, 3, 1, 1, 1, 1}, // I
	{2, 0, 1, 2, 1, 1, 1, 2, 0, 1, 2, 0, 0, 1, 1, 1}, // L
	{0, 0, 1, 2, 1, 2, 1, 1, 0, 1, 2, 2, 1, 1, 0, 1}, // J
	{1, 0, 1, 2, 1, 1, 2, 1, 0, 1, 2, 1, 1, 0, 1, 1}, // T
	{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}, // O
	{1, 2, 0, 1, 1, 1, 2, 2, 1, 2, 0, 1, 0, 0, 1, 1}, // S
	{0, 1, 1, 2, 2, 1, 2, 1, 0, 1, 1, 2, 1, 0, 1, 0}, // Z
}
var tetrominoY = [7][16]int{
	{1, 1, 1, 1, 0, 1, 2, 3, 2, 2, 2, 2, 0, 1, 2, 3}, // I
	{0, 1, 1, 1, 0, 1, 2, 2, 1, 1, 1, 2, 0, 0, 1, 2}, // L
	{0, 1, 1, 1, 0, 0, 1, 2, 1, 1, 1, 2, 0, 1, 2, 2}, // J
	{0, 1, 1, 1, 0, 1, 1, 2, 1, 1, 1, 2, 0, 1, 1, 2}, // T
	{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, // O
	{0, 0, 1, 1, 0, 1, 1, 2, 1, 1, 2, 2, 0, 1, 1, 2}, // S
	{0, 0, 1, 1, 0, 1, 1, 2, 1, 1, 2, 2, 0, 1, 1, 2}, // Z
}

type tetromino struct {
	tIndex               int
	variation            int
	positionX, positionY int
}

func (t *tetromino) TryMoveLeft() {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			if (tetrominoX[t.tIndex][i] + (t.positionX - 1)) < 0 {
				return
			}
		}
		t.positionX--
	}
}

func (t *tetromino) TryMoveRight() {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			if (tetrominoX[t.tIndex][i] + (t.positionX + 1)) >= boardSizeX {
				return
			}
		}
		t.positionX++
	}
}

func (t *tetromino) TryRotate() {
	if t != nil {
		newVariation := (t.variation + 4) % 16
		for i := newVariation; i < newVariation+4; i++ {
			newX := tetrominoX[t.tIndex][i] + t.positionX
			newY := tetrominoY[t.tIndex][i] + t.positionY
			if (newX < 0) || (newX >= boardSizeX) || (newY >= boardSizeY) {
				return
			}
		}
		t.variation = newVariation
	}
}

func (t *tetromino) TryDrop() dropResult {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			if (tetrominoY[t.tIndex][i] + (t.positionY + 1)) >= boardSizeY {
				return fix
			}
		}
		t.positionY++
	}
	return none
}

func (t *tetromino) Fix(board []int) {
	if t != nil {
		for i := 0; i < 4; i++ {
			// TODO Fix currents tetrominoes position on the board
		}
	}
}

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

func drawTmino(col termbox.Attribute, left, top int, t *tetromino) {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			tbPrint((t.positionX+tetrominoX[t.tIndex][i])*xScale+left, (t.positionY+tetrominoY[t.tIndex][i])*yScale+top, col, col, "[]")
		}
	}
}

func spawnTmino() *tetromino {
	rand.Seed(time.Now().UnixNano())
	return &tetromino{tIndex: rand.Intn(len(tetrominoX)), variation: 0, positionX: 3, positionY: -2}
}

func draw(t *tetromino) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	w, h := termbox.Size()

	centerY := h / 2
	centerX := w / 2

	boardLeft := centerX - boardSizeX*xScale/2
	boardTop := centerY - boardSizeY*yScale/2 - 2

	drawBoardFrame(coldef, boardLeft, boardTop)
	drawHints(coldef, boardLeft, boardTop)

	drawTmino(coldef, boardLeft, boardTop, t)

	termbox.Flush()
}

func onKeyDown() {
	spawnTmino()
}

func mainloop(eventQueue chan termbox.Event, done chan bool) {
	board := make([]int, 20)
	tmino := spawnTmino()
	draw(tmino)
	updateInterval := initialDropInterval
	updateTick := time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					done <- true
				case termbox.KeyArrowDown:
					updateTick.Stop()
					updateTick = time.NewTicker(time.Duration(hardDropInterval) * time.Millisecond)
				case termbox.KeyArrowLeft:
					tmino.TryMoveLeft()
				case termbox.KeyArrowRight:
					tmino.TryMoveRight()
				case termbox.KeyArrowUp:
					tmino.TryRotate()
				default:
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		case <-updateTick.C:
			switch tmino.TryDrop() {
			case fix:
				tmino.Fix(board)
				updateTick.Stop()
				updateTick = time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
				tmino = spawnTmino()
			}
		}
		draw(tmino)
	}
}

func main() {

	// initialize ui lib
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	// fetch user inputs
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// start main game loop
	done := make(chan bool)
	go mainloop(eventQueue, done)
	<-done
}
