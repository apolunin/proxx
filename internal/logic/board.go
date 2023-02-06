package logic

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	MaxBoardSize = 99
	MinBoardSize = 4
)

type (
	// board represents game field.
	board struct {
		size            int     // board size
		cells           []*cell // represents all board cells
		numVisibleCells int     // represents number of visible cells (clicked or revealed)
		numHoles        int     // represents number of holes on a board
	}

	// boardParams is used to initialize board.
	boardParams struct {
		size     int   // board size
		numHoles int   // number of holes on a board
		seed     int64 // seed for randomization, leave this hatch for testing
	}
)

// newBoard creates and initializes new game board.
func newBoard(params *boardParams) (*board, error) {
	if params.size > MaxBoardSize {
		return nil, fmt.Errorf("maximum allowed board size is %d", MaxBoardSize)
	}

	if maxHoles := params.size*params.size - 9; params.numHoles > maxHoles {
		return nil, fmt.Errorf("for size %d maximum number of holes is %d",
			params.size, params.numHoles)
	}

	seed := params.seed
	if seed == 0 {
		seed = time.Now().Unix()
	}

	var (
		random     = rand.New(rand.NewSource(seed))
		cells      = make([]*cell, params.size*params.size)
		holesSlice = getHolesIndices(params.size, params.numHoles, random)
		holesMap   = makeHolesMap(holesSlice)
		brd        = &board{
			size:            params.size,
			cells:           cells,
			numHoles:        params.numHoles,
			numVisibleCells: 0,
		}
	)

	for row := 0; row < params.size; row++ {
		for col := 0; col < params.size; col++ {
			var (
				idx         = brd.getCellIndex(row, col)
				numAdjHoles = brd.countAdjacentHoles(holesMap, row, col)
				_, isHole   = holesMap[idx]
			)

			cells[idx] = newCell(&cellParams{
				visible:     false,
				hole:        isHole,
				numAdjHoles: numAdjHoles,
				row:         row,
				col:         col,
			})
		}
	}

	return brd, nil
}

// getCell returns a pointer to cell on a board by given row and column.
func (b *board) getCell(row, col int) *cell {
	return b.cells[b.getCellIndex(row, col)]
}

// getCellIndex returns flatten (single-dimension) index of the cell.
func (b *board) getCellIndex(row, col int) int {
	return row*b.size + col
}

// getAdjacentCellIndices returns flatten indices of cells
// adjacent to cell specified by given row and column.
func (b *board) getAdjacentCellIndices(row, col int) (res []int) {
	indices := [][]int{
		{row - 1, col - 1},
		{row - 1, col},
		{row - 1, col + 1},

		{row, col - 1},
		{row, col + 1},

		{row + 1, col - 1},
		{row + 1, col},
		{row + 1, col + 1},
	}

	for _, p := range indices {
		if r, c := p[0], p[1]; (r >= 0 && r < b.size) && (c >= 0 && c < b.size) {
			res = append(res, b.getCellIndex(r, c))
		}
	}

	return
}

// countAdjacentHoles counts total holes adjacent to cell specified by given row and column.
func (b *board) countAdjacentHoles(holesMap map[int]struct{}, row, col int) (res int) {
	adjCells := b.getAdjacentCellIndices(row, col)

	for _, adjCell := range adjCells {
		if _, ok := holesMap[adjCell]; ok {
			res++
		}
	}

	return
}

// String returns a string representation of the board state combined with coordinate system.
func (b *board) String() string {
	var builder strings.Builder

	for row := -1; row < b.size; row++ {
		for col := -1; col < b.size; col++ {
			// Print column numbers
			if row < 0 {
				if col >= 0 {
					builder.WriteString(fmt.Sprintf("%3d", col+1))
				} else {
					builder.WriteString(fmt.Sprintf("%3s", ""))
				}

				continue
			}

			// Print row numbers
			if col < 0 {
				if row >= 0 {
					builder.WriteString(fmt.Sprintf("%3d", row+1))
				}

				continue
			}

			builder.WriteString(fmt.Sprintf("%3v", b.getCell(row, col)))
		}

		builder.WriteRune('\n')
	}

	return builder.String()
}

// clickCell clicks a cell specified by given row and column. If clicked cell has no adjacent holes
// then adjacent cells are clicked as well.
func (b *board) clickCell(row, col int) error {
	var click func(row, col int) error

	visitedCells := make(map[int]struct{})

	click = func(row, col int) error {
		if row < 0 || row >= b.size {
			return fmt.Errorf("row should be in range [0; %d]", b.size-1)
		}

		if col < 0 || col >= b.size {
			return fmt.Errorf("col should be in range [0; %d]", b.size-1)
		}

		var (
			clickedCellIndex = b.getCellIndex(row, col)
			clickedCell      = b.cells[clickedCellIndex]
		)

		if clickedCell.visible {
			return nil
		}

		visitedCells[clickedCellIndex] = struct{}{}
		clickedCell.click()
		b.numVisibleCells++

		if clickedCell.numAdjHoles == 0 {
			for _, adjCellIdx := range b.getAdjacentCellIndices(row, col) {
				if _, visited := visitedCells[adjCellIdx]; visited {
					continue
				}

				adjCell := b.cells[adjCellIdx]
				if err := click(adjCell.row, adjCell.col); err != nil {
					return fmt.Errorf("failed to click cell (%d, %d): %w",
						adjCell.row, adjCell.col, err)
				}

				visitedCells[adjCellIdx] = struct{}{}
			}
		}

		return nil
	}

	return click(row, col)
}

// hiddenCellsCount returns the number of unrevealed cells.
func (b *board) hiddenCellsCount() int {
	return b.size*b.size - b.numVisibleCells
}

// getHolesIndices randomly picks cells which will contain holes and
// returns their flatten (single-dimension) indices.
func getHolesIndices(size, numHoles int, r *rand.Rand) (res []int) {
	res = make([]int, size*size)

	for i := 0; i < len(res); i++ {
		res[i] = i
	}

	shuffle(res, r)

	return res[:numHoles]
}

// makeHolesMap creates a map which helps to quickly figure out if cell is a hole or not.
func makeHolesMap(holes []int) (res map[int]struct{}) {
	res = make(map[int]struct{}, len(holes))

	for _, index := range holes {
		res[index] = struct{}{}
	}

	return
}

// shuffle shuffles a slice of integers.
func shuffle[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s []T, r *rand.Rand) {
	for len(s) > 0 {
		var (
			length = len(s)
			index  = r.Intn(length)
		)

		s[length-1], s[index] = s[index], s[length-1]
		s = s[:length-1]
	}
}
