package connect4

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Turn int

const (
	Draw Turn = iota
	Player1Turn
	Player2Turn
	Player1Won
	Player2Won
)

type PlayerID int

const (
	Player1 PlayerID = 1
	Player2          = 2
)

type GameState struct {
	board         *Board
	turn          Turn
	moveListeners []chan<- Move
}

func (gameState *GameState) GetTurn() Turn {
	return gameState.turn
}

func (gameState *GameState) RegisterMoveListener(moveListener chan<- Move) {
	gameState.moveListeners = append(gameState.moveListeners, moveListener)
}

func (gameState *GameState) IsValidMove(move Move) bool {
	if gameState.IsGameOver() {
		return false
	}

	if move < 0 || int(move) >= BoardWidth {
		return false
	}

	return gameState.board[0][move] == EmptyPiece
}

func (gameState *GameState) GetPossibleMoves() []Move {
	moves := make([]Move, 0, 1)

	if gameState.IsGameOver() {
		return moves
	}

	for column := 0; column < BoardWidth; column++ {
		move := Move(column)
		if gameState.IsValidMove(move) {
			moves = append(moves, move)
		}
	}

	return moves
}

func (gameState *GameState) String() string {
	output := gameState.board.String()
	switch gameState.turn {
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

func (gameState *GameState) Print() {
	print(gameState.String())
}

func (gameState *GameState) IsGameOver() bool {
	return gameState.turn != Player1Turn && gameState.turn != Player2Turn
}

func (gameState *GameState) verifyEndGame() {
	if gameState.turn != Player1Turn && gameState.turn != Player2Turn {
		return
	}

	///////////////////////////
	// Search for 4 in a row //
	///////////////////////////

	// Search horizontally
	for y := 0; y < BoardHeight; y++ {
		player1PieceCount := 0
		player2PieceCount := 0
		x := 0
		for ; x < 3; x++ {
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}

		for ; x < BoardWidth; x++ {
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				gameState.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				gameState.turn = Player2Won
				return
			}

			switch gameState.board[y][x-3] {
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
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}

		for ; y < BoardHeight; y++ {
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				gameState.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				gameState.turn = Player2Won
				return
			}

			switch gameState.board[y-3][x] {
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
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			x++
			y++
		}

		for x < BoardWidth && y < BoardHeight {
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				gameState.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				gameState.turn = Player2Won
				return
			}

			switch gameState.board[y-3][x-3] {
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
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			x++
			y--
		}

		for x < BoardWidth && y >= 0 {
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			if player1PieceCount == 4 {
				gameState.turn = Player1Won
				return
			} else if player2PieceCount == 4 {
				gameState.turn = Player2Won
				return
			}

			switch gameState.board[y+3][x-3] {
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
		if gameState.IsValidMove(move) {
			haveValidMove = true
			break
		}
	}

	if !haveValidMove {
		gameState.turn = Draw
	}
}

func (gameState *GameState) MakeMove(move Move) error {
	if !gameState.IsValidMove(move) {
		return errors.New("Invalid Move!")
	}

	for y := BoardHeight - 1; y >= 0; y-- {
		if gameState.board[y][move] == EmptyPiece {
			switch gameState.turn {
			case Player1Turn:
				gameState.board[y][move] = Player1Piece
			case Player2Turn:
				gameState.board[y][move] = Player2Piece
			default:
				log.Fatal("Invalid Move!")
			}
			break
		}
	}

	if gameState.turn == Player1Turn {
		gameState.turn = Player2Turn
	} else if gameState.turn == Player2Turn {
		gameState.turn = Player1Turn
	}

	gameState.verifyEndGame()

	for _, moveListener := range gameState.moveListeners {
		moveListener <- move
	}

	if gameState.IsGameOver() {
		for _, moveListener := range gameState.moveListeners {
			close(moveListener)
		}
	}

	return nil
}

func NewGame() *GameState {
	return &GameState{&Board{}, Player1Turn, nil}
}

func (gameState *GameState) Clone() *GameState {
	return &GameState{gameState.board.Clone(), gameState.turn, nil}
}

/* GameState File Format:7x6 array of (RY )
 * Turn is determined by count of R vs Y
 */

func (gameState *GameState) Save(filename string) error {
	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(gameState.String())
	return nil
}

func ParseGame(gameDescription string) (*GameState, error) {
	gameState := &GameState{&Board{}, Player1Turn, nil}

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
				gameState.board[y][x] = Player1Piece
				player1PieceCount++
				expectingPiece = false
			case 'Y':
				gameState.board[y][x] = Player2Piece
				player2PieceCount++
				expectingPiece = false
			case '|':
				gameState.board[y][x] = EmptyPiece
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
		gameState.turn = Player1Turn
	} else if player1PieceCount == player2PieceCount+1 {
		gameState.turn = Player2Turn
	} else {
		gameState.Print()
		return nil, errors.New(fmt.Sprintf("invalid gameState description: (%d red pieces, %d yellow pieces)", player1PieceCount, player2PieceCount))
	}

	gameState.verifyEndGame()

	return gameState, nil
}

func LoadGame(filename string) (*GameState, error) {
	gameDescriptionBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseGame(string(gameDescriptionBytes))
}
