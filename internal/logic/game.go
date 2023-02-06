package logic

import "fmt"

type (
	// Game represents an instance of the game.
	Game struct {
		brd   *board    // represents a game board
		state GameState // represents a game state
	}

	// GameParams is used to initialize a Game instance.
	GameParams struct {
		Size     int   // tracks a game board size
		NumHoles int   // tracks number of holes on a board
		Seed     int64 // used to randomize holes position on a board
	}

	// GameState represents current game state.
	GameState int
)

const (
	InProgress GameState = iota
	Won
	Lost
)

// NewGame creates and initializes a Game instance.
func NewGame(params *GameParams) (*Game, error) {
	brd, err := newBoard(&boardParams{
		size:     params.Size,
		numHoles: params.NumHoles,
		seed:     params.Seed,
	})

	if err != nil {
		return nil, fmt.Errorf("failed tp create board: %w", err)
	}

	return &Game{
		state: InProgress,
		brd:   brd,
	}, nil
}

// ClickCell clicks the cell specified by given row and col and checks game win/loose conditions.
// Returns new game state (which becomes relevant after the click) or error if any.
func (g *Game) ClickCell(row, col int) (GameState, error) {
	if g.state != InProgress {
		return g.state, nil
	}

	if err := g.brd.clickCell(row, col); err != nil {
		return g.state, err
	}

	if clickedCell := g.brd.getCell(row, col); clickedCell.hole {
		g.state = Lost
		return g.state, nil
	}

	if g.brd.hiddenCellsCount() == g.brd.numHoles {
		g.state = Won
		return g.state, nil
	}

	return g.state, nil
}

// String returns a string representation of the game.
func (g *Game) String() string {
	return g.brd.String()
}
