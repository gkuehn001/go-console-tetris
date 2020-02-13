package main

import "math/rand"

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

func (t *tetromino) TryMoveLeft(board []int) {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			x := tetrominoX[t.tIndex][i] + t.positionX
			y := tetrominoY[t.tIndex][i] + t.positionY
			if !checkFreeBoardPosition(board, x-1, y) {
				return
			}
		}
		t.positionX--
	}
}

func (t *tetromino) TryMoveRight(board []int) {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			x := tetrominoX[t.tIndex][i] + t.positionX
			y := tetrominoY[t.tIndex][i] + t.positionY
			if !checkFreeBoardPosition(board, x+1, y) {
				return
			}
		}
		t.positionX++
	}
}

func (t *tetromino) TryRotate() {
	if t != nil {
		index := t.tIndex
		newVariation := (t.variation + 4) % len(tetrominoX[t.tIndex])
		if glGameMode == gmChaos {
			index = rand.Intn(len(tetrominoX))
			newVariation = rand.Intn(4) * 4
		}
		for i := newVariation; i < newVariation+4; i++ {
			newX := tetrominoX[index][i] + t.positionX
			newY := tetrominoY[index][i] + t.positionY
			if (newX < 0) || (newX >= boardSizeX) || (newY >= boardSizeY) {
				return
			}
		}
		t.tIndex = index
		t.variation = newVariation
	}
}

func (t *tetromino) TryDrop(board []int) dropResult {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			x := tetrominoX[t.tIndex][i] + t.positionX
			y := tetrominoY[t.tIndex][i] + t.positionY
			if !checkFreeBoardPosition(board, x, y+1) {
				return fix
			}
		}
		t.positionY++
	}
	return none
}

func (t *tetromino) Fix(board []int) bool {
	if t != nil {
		for i := t.variation; i < t.variation+4; i++ {
			x := tetrominoX[t.tIndex][i] + t.positionX
			y := tetrominoY[t.tIndex][i] + t.positionY
			if y < 0 {
				return false
			}
			board[y] = board[y] | (1 << x)
		}
	}
	return true
}
