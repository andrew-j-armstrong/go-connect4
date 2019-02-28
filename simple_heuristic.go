package main

type SimpleHeuristic struct {
	targetPlayer Connect4Player
}

func NewSimpleHeuristic(targetPlayer Connect4Player) *SimpleHeuristic {
	return &SimpleHeuristic{targetPlayer}
}

func (heuristic *SimpleHeuristic) Heuristic(game *Game) float64 {
	if !game.IsGameOver() || game.turn == Draw {
		return 0.0
	} else if game.turn == Player1Won {
		if heuristic.targetPlayer == Player1 {
			return 1.0
		} else {
			return -1.0
		}
	} else {
		if heuristic.targetPlayer == Player1 {
			return -1.0
		} else {
			return 1.0
		}
	}
}
