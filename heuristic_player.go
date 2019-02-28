package main

type HeuristicPlayerHeuristic interface {
	Heuristic(game *Game) float64
}

type HeuristicPlayer struct {
	game      *Game
	heuristic HeuristicPlayerHeuristic
}

func (player *HeuristicPlayer) GetNextMove() Move {
	getMoveHeuristic := func(move interface{}) float64 {
		game := player.game.Clone()
		game.MakeMoveGeneric(move)
		return player.heuristic.Heuristic(game)
	}
	return player.game.GetPossibleMovesGeneric().GetBestEntry(getMoveHeuristic).(Move)
}

func NewHeuristicPlayer(game *Game, heuristic HeuristicPlayerHeuristic) *HeuristicPlayer {
	return &HeuristicPlayer{game, heuristic}
}
