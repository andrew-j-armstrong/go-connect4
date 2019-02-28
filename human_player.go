package main

import (
	"fmt"
)

type HumanPlayer struct {
	game *Game
}

func (player *HumanPlayer) GetNextMove() Move {
	player.game.Print()

	var move Move
	for {
		print("Column: ")

		var column int
		_, err := fmt.Scan(&column)
		if err != nil {
			fmt.Println(err)
			continue
		}

		move = Move(column)

		if !player.game.IsValidMove(move) {
			fmt.Println("Invalid Move!")
			continue
		}

		break
	}

	return move
}

func NewHumanPlayer(game *Game) *HumanPlayer {
	return &HumanPlayer{game}
}
