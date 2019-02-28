package main

import (
	"fmt"
	"log"
	"strconv"
)

type Move int

func ParseMove(moveString string) Move {
	move, err := strconv.Atoi(moveString)
	if err != nil {
		log.Fatalf("unable to parse move: %s", err)
	}

	if move < 0 || int(move) >= BoardWidth {
		log.Fatalf("invalid move: %d", move)
	}

	return Move(move)
}

func (move Move) String() string {
	return fmt.Sprintf("%d", move)
}
