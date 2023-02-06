package logic

import (
	"fmt"
	"testing"
)

const (
	size     = 4
	numHoles = 5
	seed     = 1
)

func TestGameWon(t *testing.T) {
	// Holes location in test
	//      1  2  3  4
	//   1  0  1  3  ◯
	//   2  1  3  ◯  ◯
	//   3  1  ◯  ◯  3
	//   4  1  2  2  1

	game, err := NewGame(&GameParams{
		Size:     size,
		NumHoles: numHoles,
		Seed:     seed,
	})

	if err != nil {
		t.Fatalf("cannot create game: %v", err)
	}

	steps := [][]int{
		{1, 1}, {1, 2}, {1, 3},
		{2, 1}, {2, 2},
		{3, 1}, {3, 4},
		{4, 1}, {4, 2}, {4, 3}, {4, 4},
	}

	gameState := InProgress

	for _, pair := range steps {
		row, col := pair[0]-1, pair[1]-1

		if gameState, err = game.ClickCell(row, col); err != nil {
			t.Fatalf("failed to click cell (%d,%d): %v", row, col, err)
		}
	}

	fmt.Printf("\n%v\n", game)

	if gameState != Won {
		t.Errorf("expected game state %d, got %d", Won, gameState)
	}
}

func TestGameLost(t *testing.T) {
	// Holes location in test
	//      1  2  3  4
	//   1  0  1  3  ◯
	//   2  1  3  ◯  ◯
	//   3  1  ◯  ◯  3
	//   4  1  2  2  1

	game, err := NewGame(&GameParams{
		Size:     size,
		NumHoles: numHoles,
		Seed:     seed,
	})

	if err != nil {
		t.Fatalf("cannot create game: %v", err)
	}

	steps := [][]int{
		{1, 1}, {1, 2}, {1, 3},
		{2, 1}, {2, 2},
		{3, 1}, {3, 3}, {3, 4},
	}

	gameState := InProgress

	for _, pair := range steps {
		row, col := pair[0]-1, pair[1]-1

		if gameState, err = game.ClickCell(row, col); err != nil {
			t.Fatalf("failed to click cell (%d,%d): %v", row, col, err)
		}
	}

	fmt.Printf("\n%v\n", game)

	if gameState != Lost {
		t.Errorf("expected game state %d, got %d", Won, gameState)
	}
}
