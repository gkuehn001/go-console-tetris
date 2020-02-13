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
const fullRow = (1 << boardSizeX) - 1

type dropResult int

const (
	none dropResult = 0
	fix  dropResult = 1
)

type gameMode int

const (
	gmNormal gameMode = 0
	gmChaos  gameMode = 1
	gmMax    gameMode = 2
)

var glGameMode gameMode = gmNormal

func checkFreeBoardPosition(board []int, x, y int) bool {
	if x < 0 || x >= boardSizeX || y >= boardSizeY {
		return false
	}
	if y >= 0 && (board[y]&(1<<x)) == (1<<x) {
		return false
	}
	return true
}

func removeFullRows(board []int) []int {
	for i := boardSizeY - 1; i >= 0; i-- {
		if board[i] == fullRow {
			board = append([]int{0}, append(board[:i], board[i+1:]...)...)
			i++
		}
	}
	return board
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func drawHints(col termbox.Attribute, left, top int) {
	x, y := left-24, top+2
	tbPrint(x, y, col, col, "     \u2190 : Move Left")
	y++
	tbPrint(x, y, col, col, "     \u2192 : Move Right")
	y++
	switch glGameMode {
	case gmNormal:
		tbPrint(x, y, col, col, "     \u2191 : Rotate")
	case gmChaos:
		tbPrint(x, y, col, col, "     \u2191 : Switch")
	}
	y++
	tbPrint(x, y, col, col, "     \u2193 : Drop")
	y += 2
	switch glGameMode {
	case gmNormal:
		tbPrint(x, y, col, col, "CTRL-M : Normal")
	case gmChaos:
		tbPrint(x, y, col, col, "CTRL-M : Chaos")
	}

	y += 2
	tbPrint(x, y, col, col, " SPACE : Restart")
	y++
	tbPrint(x, y, col, col, "   ESC : Quit")

	tbPrint(left+4, top+boardSizeY*yScale+3, col, col, "Play Tetris!")
}

func drawGameOver(col termbox.Attribute, centerX int, centerY int) {
	tbPrint(centerX-9, centerY-4, col, col, "+----------------+")
	tbPrint(centerX-9, centerY-3, col, col, "|   Game Over!   |")
	tbPrint(centerX-9, centerY-2, col, col, "+----------------+")
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
	return &tetromino{tIndex: rand.Intn(len(tetrominoX)), variation: 0, positionX: boardSizeX/2 - 2, positionY: -2}
}

func drawBoard(col termbox.Attribute, left, top int, board []int) {
	for y, by := range board {
		if by > 0 {
			for x := 0; x < boardSizeX; x++ {
				if (by & (1 << x)) == (1 << x) {
					tbPrint(left+x*xScale, top+y*yScale, col, col, "[]")
				}
			}
		}
	}
}

func draw(t *tetromino, board []int, gameOver bool) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	w, h := termbox.Size()

	centerY := h / 2
	centerX := w / 2

	boardLeft := centerX - boardSizeX*xScale/2
	boardTop := centerY - boardSizeY*yScale/2 - 2

	drawBoardFrame(coldef, boardLeft, boardTop)
	drawHints(coldef, boardLeft, boardTop)

	drawBoard(coldef, boardLeft, boardTop, board)

	drawTmino(coldef, boardLeft, boardTop, t)

	if gameOver {
		drawGameOver(coldef, centerX, centerY)
	}

	termbox.Flush()
}

func mainloop(eventQueue chan termbox.Event, done chan bool) {
	gameOver := false
	board := make([]int, boardSizeY)
	tmino := spawnTmino()
	draw(tmino, board, gameOver)
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
					if !gameOver {
						updateTick.Stop()
						updateTick = time.NewTicker(time.Duration(hardDropInterval) * time.Millisecond)
					}
				case termbox.KeyArrowLeft:
					if !gameOver {
						tmino.TryMoveLeft(board)
					}
				case termbox.KeyArrowRight:
					if !gameOver {
						tmino.TryMoveRight(board)
					}
				case termbox.KeyArrowUp:
					if !gameOver {
						tmino.TryRotate()
					}
				case termbox.KeySpace:
					gameOver = false
					updateTick = time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
					board = make([]int, boardSizeY)
					tmino = spawnTmino()
				case termbox.KeyCtrlM:
					glGameMode = (glGameMode + 1) % gmMax
				default:
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		case <-updateTick.C:
			switch tmino.TryDrop(board) {
			case fix:
				if !gameOver {
					if !tmino.Fix(board) {
						// Game over
						updateTick.Stop()
						gameOver = true
					} else {
						board = removeFullRows(board)
						updateTick.Stop()
						updateTick = time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
						tmino = spawnTmino()
					}
				}
			}
		}
		draw(tmino, board, gameOver)
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

	rand.Seed(time.Now().UnixNano())

	// start main game loop
	done := make(chan bool)
	go mainloop(eventQueue, done)
	<-done
}
