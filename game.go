package chesscli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/games"
)

// ...
var (
	ErrExit = errors.New("exit")
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
	game        GameModel
	stringer    Stringer
	reader      *bufio.Reader
	writer      io.Writer
	exitCommand string
}

// NewGame ...
func NewGame(
	game GameModel,
	stringer Stringer,
	reader io.Reader,
	writer io.Writer,
	exitCommand string,
) Game {
	return Game{
		game:        game,
		stringer:    stringer,
		reader:      bufio.NewReader(reader),
		writer:      writer,
		exitCommand: exitCommand,
	}
}

// WritePrompt ...
func (game Game) WritePrompt(
	prompt string,
) error {
	text := game.stringer(game.game.Storage())
	text += "\n"

	_, err :=
		io.WriteString(game.writer, text)
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

	prompt += " "

	_, err =
		io.WriteString(game.writer, prompt)
	if err != nil {
		const message = "unable to write " +
			"the prompt message: %s"
		return fmt.Errorf(message, err)
	}

	return nil
}

// ReadMove ...
func (game Game) ReadMove(
	prompt string,
) error {
	err := game.WritePrompt(prompt)
	switch err {
	case nil:
	case games.ErrCheckmate, games.ErrDraw:
		return err // don't wrap
	default:
		const message = "unable to write " +
			"the prompt: %s"
		return fmt.Errorf(message, err)
	}

	command, err :=
		game.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		const message = "unable to read " +
			"the command: %s"
		return fmt.Errorf(message, err)
	}

	command = strings.TrimSpace(command)
	command = strings.ToLower(command)
	if command == game.exitCommand {
		return ErrExit // don't wrap
	}

	move, err := uci.DecodeMove(command)
	if err != nil {
		const message = "unable to decode " +
			"the move: %s"
		return fmt.Errorf(message, err)
	}

	err = game.game.ApplyMove(move)
	if err != nil {
		const message = "unable to apply " +
			"the move: %s"
		return fmt.Errorf(message, err)
	}

	return nil
}

// SearchMove ...
func (game Game) SearchMove(
	prompt string,
) error {
	err := game.WritePrompt(prompt)
	switch err {
	case nil:
	case games.ErrCheckmate, games.ErrDraw:
		return err // don't wrap
	default:
		const message = "unable to write " +
			"the prompt: %s"
		return fmt.Errorf(message, err)
	}

	move := game.game.SearchMove()
	fmt.Println(uci.EncodeMove(move))

	return nil
}
