package main

type Piece int

const (
	EmptyPiece Piece = iota
	Player1Piece
	Player2Piece
)

const BoardHeight int = 6
const BoardWidth int = 7

type Board [BoardHeight][BoardWidth]Piece

func (board *Board) IsEqual(otherBoard *Board) bool {
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board[y][x] != otherBoard[y][x] {
				return false
			}
		}
	}

	return true
}

func (board *Board) String() string {
	var output string
	for y := 0; y < BoardHeight; y++ {
		output += "+---+---+---+---+---+---+---+\n"
		for x := 0; x < BoardWidth; x++ {
			output += "| "
			switch board[y][x] {
			case EmptyPiece:
				output += "  "
			case Player1Piece:
				output += "R "
			case Player2Piece:
				output += "Y "
			default:
				output += "? "
			}
		}
		output += "|\n"
	}
	output += "+---+---+---+---+---+---+---+\n"

	return output
}

func (board *Board) Print() {
	print(board.String())
}

func (board *Board) Clone() *Board {
	newBoard := &Board{}

	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			newBoard[y][x] = board[y][x]
		}
	}

	return newBoard
}

func (board *Board) CloneGeneric() interface{} {
	return board.Clone()
}
