package connect4

import (
	"fmt"
	"strconv"
)

type Move int

func ParseMove(moveString string) (Move, error) {
	move, err := strconv.Atoi(moveString)
	if err != nil {
		return 0, fmt.Errorf("unable to parse move: %s", err)
	}

	if move < 0 || int(move) >= BoardWidth {
		return 0, fmt.Errorf("invalid move: %d", move)
	}

	return Move(move), nil
}

func (move Move) String() string {
	return fmt.Sprintf("%d", move)
}
