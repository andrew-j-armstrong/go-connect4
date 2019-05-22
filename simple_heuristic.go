package connect4

type SimpleHeuristic struct {
	targetPlayer PlayerID
}

func NewSimpleHeuristic(targetPlayer PlayerID) *SimpleHeuristic {
	return &SimpleHeuristic{targetPlayer}
}

func (heuristic *SimpleHeuristic) Heuristic(gameState *GameState) float64 {
	if !gameState.IsGameOver() || gameState.turn == Draw {
		return 0.0
	} else if gameState.turn == Player1Won {
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
