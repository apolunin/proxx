package logic

import "strconv"

const (
	square = "□"
	circle = "◯"
)

type (
	// cell represents a game board cell.
	cell struct {
		hole        bool // tracks if this cell is a hole
		visible     bool // tracks if this cell is visible
		numAdjHoles int  // tracks the number of adjacent holes
		row, col    int  // tracks row and column of this cell on a board
	}

	// cellParams is used to initialize a cell.
	cellParams struct {
		hole        bool // tracks if this cell is a hole
		visible     bool // tracks if this cell is visible
		numAdjHoles int  // tracks the number of adjacent holes
		row, col    int  // tracks row and column of this cell on a board
	}
)

// newCell initializes a cell and returns a pointer to it.
func newCell(params *cellParams) *cell {
	return &cell{
		hole:        params.hole,
		visible:     params.visible,
		numAdjHoles: params.numAdjHoles,
		row:         params.row,
		col:         params.col,
	}
}

// click clicks the cell (i.e. reveals it on a board).
func (c *cell) click() {
	c.visible = true
}

// String returns a string representation of the cell on a board.
func (c *cell) String() string {
	if !c.visible {
		return square
	}

	if c.hole {
		return circle
	}

	return strconv.Itoa(c.numAdjHoles)
}
