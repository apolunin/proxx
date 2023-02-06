package main

import (
	"flag"
	"fmt"
	"time"

	"proxx/internal/logic"
)

func main() {
	const (
		minHoles = 5
		maxHoles = logic.MaxBoardSize*logic.MaxBoardSize - 9
	)

	var (
		size = flag.Int("n", 4, fmt.Sprintf(
			"specifies game field size (should be in range [%d;%d])",
			logic.MinBoardSize,
			logic.MaxBoardSize,
		))

		numHoles = flag.Int("k", 5, fmt.Sprintf(
			"specifies number of holes (should be in range [%d;%d])",
			minHoles,
			maxHoles,
		))
	)

	flag.Parse()

	game, err := logic.NewGame(&logic.GameParams{
		Size:     *size,
		NumHoles: *numHoles,
		Seed:     time.Now().Unix(),
	})

	if err != nil {
		return
	}

	gameState := logic.InProgress

	for gameState == logic.InProgress {
		fmt.Printf("\n%v\n", game)

		row, col, err := readCellCoords(*size)
		if err != nil {
			fmt.Printf("failed to read user input: %v\naborting the game...\n",
				err)
			return
		}

		if gameState, err = game.ClickCell(row, col); err != nil {
			fmt.Printf("failed to click cell (%d, %d): %v\naborting the game...\n",
				row, col, err)
			return
		}
	}

	fmt.Printf("\n%v\n", game)

	switch gameState {
	case logic.Won:
		fmt.Println("*** YOU WON ***")
	case logic.Lost:
		fmt.Println("*** YOU LOST ***")
	}
}

func readCellCoords(size int) (row, col int, err error) {
	for {
		fmt.Print("row, col = ")

		if _, err = fmt.Scanf("%d,%d", &row, &col); err != nil {
			fmt.Printf("input error: %v\n", err)
			continue
		}

		if row >= 1 && row <= size && col >= 1 && col <= size {
			break
		}

		fmt.Printf("row and col should be in range [1; %d]\n", size)
	}

	row--
	col--

	return
}
