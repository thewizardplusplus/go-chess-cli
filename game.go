package chesscli

import (
	"fmt"
	"io"

	models "github.com/thewizardplusplus/go-chess-models"
)

// GameModel ...
type GameModel interface {
	Storage() models.PieceStorage
	State() error
	ApplyMove(move models.Move) error
	SearchMove() models.Move
}

// Stringer ...
type Stringer func(
	storage models.PieceStorage,
) string

// Game ...
type Game struct {
	game     GameModel
	stringer Stringer
	writer   io.Writer
}

// NewGame ...
func NewGame(
	game GameModel,
	stringer Stringer,
	writer io.Writer,
) Game {
	return Game{
		game:     game,
		stringer: stringer,
		writer:   writer,
	}
}

// WritePrompt ...
func (game Game) WritePrompt(
	prompt string,
) error {
	text := game.stringer(game.game.Storage())
	text += "\n"

	_, err := io.WriteString(game.writer, text)
	if err != nil {
		const message = "unable to write " +
			"the storage: %s"
		return fmt.Errorf(message, err)
	}

	// it should be checked
	// after writing the storage
	state := game.game.State()
	if state != nil {
		return state // don't wrap
	}

	_, err =
		io.WriteString(game.writer, prompt)
	if err != nil {
		const message = "unable to write " +
			"the prompt message: %s"
		return fmt.Errorf(message, err)
	}

	return nil
}
