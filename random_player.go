package main

type RandomPlayer struct {
	game *Game
}

func (player *RandomPlayer) GetNextMove() Move {
	return player.game.GetPossibleMovesGeneric().SelectRandom().(Move)
}

func NewRandomPlayer(game *Game) *RandomPlayer {
	return &RandomPlayer{game}
}
