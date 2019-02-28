package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/carbon-12/go-extensions"
)

type Turn int

const (
	Draw Turn = iota
	Player1Turn
	Player2Turn
	Player1Won
	Player2Won
)

type Connect4Player int

const (
	Player1 Connect4Player = 1
	Player2                = 2
)

type Game struct {
	board                *Board
	turn                 Turn
	moveListeners        []chan<- Move
	genericMoveListeners []chan<- interface{}
}

func (game *Game) RegisterMoveListener(moveListener chan<- Move) {
	game.moveListeners = append(game.moveListeners, moveListener)
}

func (game *Game) RegisterMoveListenerGeneric(moveListener chan<- interface{}) {
	game.genericMoveListeners = append(game.genericMoveListeners, moveListener)
}

func (game *Game) IsValidMove(move Move) bool {
	if game.IsGameOver() {
		return false
	}

	if move < 0 || int(move) >= BoardWidth {
		return false
	}

	return game.board[0][move] == EmptyPiece
}

func (game *Game) IsValidMoveGeneric(move interface{}) bool {
	switch m := move.(type) {
	case Move:
		return game.IsValidMove(m)
	default:
		return false
	}
}

func (game *Game) GetPossibleMoves() []Move {
	moves := make([]Move, 0, 1)

	if game.IsGameOver() {
		return moves
	}

	for column := 0; column < BoardWidth; column++ {
		move := Move(column)
		if game.IsValidMove(move) {
			moves = append(moves, move)
		}
	}

	return moves
}

func (game *Game) GetPossibleMovesGeneric() *extensions.InterfaceSlice {
	moves := make(extensions.InterfaceSlice, 0, 1)

	if game.IsGameOver() {
		return &moves
	}

	for column := 0; column < BoardWidth; column++ {
		move := Move(column)
		if game.IsValidMove(move) {
			moves = append(moves, move)
		}
	}

	return &moves
}

func (game *Game) String() string {
	output := game.board.String()
	switch game.turn {
	case Draw:
		output += "Game Over - Draw!\n"
	case Player1Turn:
		output += "Player 1's turn.\n"
	case Player2Turn:
		output += "Player 2's turn.\n"
	case Player1Won:
		output += "Game Over - Player 1 Won!\n"
	case Player2Won:
		output += "Game Over - Player 2 Won!\n"
	default:
		output += "Invalid Turn!\n"
	}
	return output
}

func (game *Game) Print() {
	print(game.String())
}

func (game *Game) IsGameOver() bool {
	return game.turn != Player1Turn && game.turn != Player2Turn
}

func (game *Game) verifyEndGame() {
	if game.turn != Player1Turn && game.turn != Player2Turn {
		return
	}

	////////////////////////
	// Search for 4 in a row
	////////////////////////

	// Search horizontally
	for y := 0; y < BoardHeight; y++ {
		player1PieceCount := 0
		player2PieceCount := 0
		x := 0
		for ; x < 3; x++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}

		for ; x < BoardWidth; x++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				game.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				game.turn = Player2Won
				return
			}

			switch game.board[y][x-3] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}
		}
	}

	// Search vertically
	for x := 0; x < BoardWidth; x++ {
		player1PieceCount := 0
		player2PieceCount := 0
		y := 0
		for ; y < 3; y++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}

		for ; y < BoardHeight; y++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				game.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				game.turn = Player2Won
				return
			}

			switch game.board[y-3][x] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}
		}
	}

	// Search diagonally up
	for xIndex := 4 - BoardHeight; xIndex < BoardWidth-2; xIndex++ {
		player1PieceCount := 0
		player2PieceCount := 0

		var (
			x int
			y int
		)
		if xIndex < 0 {
			x = 0
			y = -xIndex
		} else {
			x = xIndex
			y = 0
		}

		for i := 0; i < 3; i++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			x++
			y++
		}

		for x < BoardWidth && y < BoardHeight {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				game.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				game.turn = Player2Won
				return
			}

			switch game.board[y-3][x-3] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}

			x++
			y++
		}
	}

	// Search diagonally down
	for xIndex := 4 - BoardHeight; xIndex < BoardWidth-3; xIndex++ {
		player1PieceCount := 0
		player2PieceCount := 0

		var (
			x int
			y int
		)
		if xIndex < 0 {
			x = 0
			y = BoardHeight - 1 + xIndex
		} else {
			x = xIndex
			y = BoardHeight - 1
		}

		for i := 0; i < 3; i++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			x++
			y--
		}

		for x < BoardWidth && y >= 0 {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				game.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				game.turn = Player2Won
				return
			}

			switch game.board[y+3][x-3] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}

			x++
			y--
		}
	}

	///////////////////////////////////////////
	// If nobody has won, then check for a draw
	///////////////////////////////////////////

	haveValidMove := false
	for x := 0; x < BoardWidth; x++ {
		move := Move(x)
		if game.IsValidMove(move) {
			haveValidMove = true
			break
		}
	}

	if !haveValidMove {
		game.turn = Draw
	}
}

func (game *Game) MakeMoveGeneric(move interface{}) error {
	switch m := move.(type) {
	case Move:
		return game.MakeMove(m)
	default:
		return errors.New("Invalid Move Type!")
	}
}

func (game *Game) MakeMove(move Move) error {
	if !game.IsValidMove(move) {
		return errors.New("Invalid Move!")
	}

	for y := BoardHeight - 1; y >= 0; y-- {
		if game.board[y][move] == EmptyPiece {
			switch game.turn {
			case Player1Turn:
				game.board[y][move] = Player1Piece
			case Player2Turn:
				game.board[y][move] = Player2Piece
			default:
				log.Fatal("Invalid Move!")
			}
			break
		}
	}

	if game.turn == Player1Turn {
		game.turn = Player2Turn
	} else if game.turn == Player2Turn {
		game.turn = Player1Turn
	}

	game.verifyEndGame()

	for _, moveListener := range game.moveListeners {
		moveListener <- move
	}

	for _, moveListener := range game.genericMoveListeners {
		moveListener <- move
	}

	if game.IsGameOver() {
		for _, moveListener := range game.moveListeners {
			close(moveListener)
		}

		for _, moveListener := range game.genericMoveListeners {
			close(moveListener)
		}
	}

	return nil
}

func NewGame() *Game {
	return &Game{&Board{}, Player1Turn, nil, nil}
}

func (game *Game) Clone() *Game {
	return &Game{game.board.Clone(), game.turn, nil, nil}
}

func (game *Game) CloneGeneric() interface{} {
	return game.Clone()
}

/* Game File Format:7x6 array of (RY )
 * Turn is determined by count of R vs Y
 */

func (game *Game) Save(filename string) {
	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error opening \"%s\" for writing: %s", filename, err)
	}

	defer f.Close()

	f.WriteString(game.String())
}

func ParseGame(gameDescription string) (*Game, error) {
	game := &Game{&Board{}, Player1Turn, nil, nil}

	player1PieceCount := 0
	player2PieceCount := 0

	y := 0
	x := 0
	expectingPiece := false
	for _, c := range gameDescription {
		if !expectingPiece {
			if c == '|' {
				expectingPiece = true
			}
		} else {
			switch c {
			case 'R':
				game.board[y][x] = Player1Piece
				player1PieceCount++
				expectingPiece = false
			case 'Y':
				game.board[y][x] = Player2Piece
				player2PieceCount++
				expectingPiece = false
			case '|':
				game.board[y][x] = EmptyPiece
			case '\n':
				expectingPiece = false
				continue
			default:
				continue
			}

			if x < BoardWidth-1 {
				x++
			} else {
				y++
				x = 0
			}

			if y >= BoardHeight {
				break
			}
		}
	}

	if player1PieceCount == player2PieceCount {
		game.turn = Player1Turn
	} else if player1PieceCount == player2PieceCount+1 {
		game.turn = Player2Turn
	} else {
		game.Print()
		return nil, errors.New(fmt.Sprintf("invalid game description: (%d red pieces, %d yellow pieces)", player1PieceCount, player2PieceCount))
	}

	game.verifyEndGame()

	return game, nil
}

func LoadGame(filename string) (*Game, error) {
	gameDescriptionBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseGame(string(gameDescriptionBytes))
}
