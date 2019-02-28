package main

import "log"

type GameSimPlayer interface {
	GetNextMove() Move
}

type GameSim struct {
	game    *Game
	player1 GameSimPlayer
	player2 GameSimPlayer
}

func (gameSim *GameSim) Run() {
	for {
		if gameSim.game.IsGameOver() {
			break
		}

		var move Move
		switch gameSim.game.turn {
		case 1:
			move = gameSim.player1.GetNextMove()
		case 2:
			move = gameSim.player2.GetNextMove()
		default:
			log.Fatal("Invalid Turn!")
		}

		err := gameSim.game.MakeMove(move)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewGameSim(game *Game, player1 GameSimPlayer, player2 GameSimPlayer) *GameSim {
	return &GameSim{game, player1, player2}
}
