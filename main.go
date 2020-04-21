package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type cell struct {
	X     float64
	Y     float64
	alive int
}

func main() {
	pixelgl.Run(run)
}

func run() {

	/*
	* -----------------Game of life setup-----------------
	 */
	const CellSize = 14

	// Each cell have 3 propertys X, Y and alive
	// Array of Cells (50 x 50)
	board := [51][51]cell{}

	// Fill the array
	for i := 0; i < len(board)-1; i++ {
		for j := 0; j < len(board[i])-1; j++ {
			board[i][j] = cell{float64(j) * CellSize, float64(i) * CellSize, 0}
		}
	}

	/*
	* -------------------Pixel setup-------------------------
	 */
	// Window config
	cfg := pixelgl.WindowConfig{
		Title:  "Press space for pause / unpause",
		Bounds: pixel.R(0, 0, 700, 700),
	}
	// Window initialization
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	p := true

	for !win.Closed() {
		//Pause unpause
		if win.JustPressed(pixelgl.KeySpace) {
			p = !p
		}

		// On click change state
		if win.JustPressed(pixelgl.MouseButton1) {
			pos := win.MousePosition()
			celX := int(pos.X / 14)
			celY := int(pos.Y / 14)
			if board[celY][celX].alive == 1 {
				board[celY][celX].alive = 0
			} else {
				board[celY][celX].alive = 1
			}
		}

		win.Clear(pixel.RGB(0.2, 0.2, 0.2))
		updateBoard(board, win)
		if !p {
			board = getNewBoard(board)
		}
		time.Sleep(50 * time.Millisecond)
		win.Update()
	}
}

func updateBoard(board [51][51]cell, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	for i := 0; i < len(board)-1; i++ {
		for j := 0; j < len(board[i])-1; j++ {
			c := board[i][j]

			// Draw the Cell whit gray if dead and white if alive
			if c.alive == 0 {
				imd.Color = pixel.RGB(0.6, 0.6, 0.6)
			} else {
				imd.Color = pixel.RGB(1, 1, 1)
			}
			imd.Push(pixel.V(c.X, c.Y))
			imd.Push(pixel.V(c.X+14, c.Y+14))
			if c.alive == 0 {
				imd.Rectangle(1)
			} else {
				imd.Rectangle(0)
			}
		}
	}
	imd.Draw(win)
}

func getNewBoard(board [51][51]cell) [51][51]cell {
	newBoard := board

	iSave := 0
	jSave := 0

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {

			iSave = i
			jSave = j

			neighbors := 0
			// Get number of neighbors of each Cell
			neighbors += board[torus(i-1)][torus(j-1)].alive
			neighbors += board[torus(i-1)][j].alive
			neighbors += board[torus(i-1)][torus(j+1)].alive // Lower
			neighbors += board[i][torus(j-1)].alive
			neighbors += board[i][torus(j+1)].alive // Same line
			neighbors += board[torus(i+1)][torus(j-1)].alive
			neighbors += board[torus(i+1)][j].alive
			neighbors += board[torus(i+1)][torus(j+1)].alive // Upper

			i = iSave
			j = jSave

			// Apply rules
			// Rule 1: A dead Cell whit 3 neighbors revive
			if board[i][j].alive == 0 && neighbors == 3 {
				newBoard[i][j].alive = 1
			}
			// Rule 2: An alive Cell dead if have less than 2 or more than 3 neighbors
			if board[i][j].alive == 1 && (neighbors < 2 || neighbors > 3) {
				newBoard[i][j].alive = 0
			}
		}
	}

	return newBoard
}

func torus(n int) int {
	if n == -1 {
		return 50
	} else if n == 51 {
		return 0
	}

	return n
}
