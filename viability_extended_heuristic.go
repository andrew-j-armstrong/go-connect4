package connect4

type ViabilityExtendedHeuristic struct {
	targetPlayer PlayerID
}

func NewViabilityExtendedHeuristic(targetPlayer PlayerID) *ViabilityExtendedHeuristic {
	return &ViabilityExtendedHeuristic{targetPlayer}
}

func (heuristic *ViabilityExtendedHeuristic) increaseViabilityScores(player1PieceCount int, player2PieceCount int, player1Viability *int, player2Viability *int) {
	if player2PieceCount == 0 {
		switch player1PieceCount {
		case 1:
			*player1Viability += 1
		case 2:
			*player1Viability += 5
		case 3:
			*player1Viability += 20
		}
	} else if player1PieceCount == 0 {
		switch player2PieceCount {
		case 1:
			*player2Viability += 1
		case 2:
			*player2Viability += 5
		case 3:
			*player2Viability += 20
		}
	}
}

func (heuristic *ViabilityExtendedHeuristic) Heuristic(gameState *GameState) float64 {
	if gameState.turn == Draw {
		return 0.0
	} else if gameState.turn == Player1Won {
		if heuristic.targetPlayer == Player1 {
			return 1.0
		} else {
			return -1.0
		}
	} else if gameState.turn == Player2Won {
		if heuristic.targetPlayer == Player1 {
			return -1.0
		} else {
			return 1.0
		}
	}

	var player1Viability int
	var player2Viability int

	// Look for next turn win opportunity
	var currentPlayerPiece Piece
	if gameState.turn == Player1Turn {
		currentPlayerPiece = Player1Piece
	} else if gameState.turn == Player2Turn {
		currentPlayerPiece = Player2Piece
	}

	currentPlayerWinOpportunities := 0
	for x := 0; x < BoardWidth; x++ {
		y := BoardHeight - 1
		for y >= 0 && gameState.board[y][x] != EmptyPiece {
			y--
		}

		if y < 0 {
			continue
		}

		var currentPlayerPieceCount int

		// Check horizontal
		currentPlayerPieceCount = 0
		for nextX := x - 1; nextX >= 0; nextX-- {
			if gameState.board[y][nextX] != currentPlayerPiece {
				break
			}
			currentPlayerPieceCount++
		}
		for nextX := x + 1; nextX < BoardWidth; nextX++ {
			if gameState.board[y][nextX] != currentPlayerPiece {
				break
			}
			currentPlayerPieceCount++
		}

		if currentPlayerPieceCount >= 3 {
			currentPlayerWinOpportunities++
			continue
		}

		// Check vertical
		if y < 3 {
			currentPlayerPieceCount = 0

			for nextY := y + 1; nextY < BoardHeight; nextY++ {
				if gameState.board[nextY][x] != currentPlayerPiece {
					break
				}
				currentPlayerPieceCount++
			}

			if currentPlayerPieceCount >= 3 {
				currentPlayerWinOpportunities++
				continue
			}
		}

		// Check diagonal up
		currentPlayerPieceCount = 0

		for nextX, nextY := x-1, y-1; nextX >= 0 && nextY >= 0; nextX, nextY = nextX-1, nextY-1 {
			if gameState.board[nextY][nextX] != currentPlayerPiece {
				break
			}
			currentPlayerPieceCount++
		}

		for nextX, nextY := x+1, y+1; nextX < BoardWidth && nextY < BoardHeight; nextX, nextY = nextX+1, nextY+1 {
			if gameState.board[nextY][nextX] != currentPlayerPiece {
				break
			}
			currentPlayerPieceCount++
		}

		if currentPlayerPieceCount >= 3 {
			currentPlayerWinOpportunities++
			continue
		}

		// Check diagonal down
		currentPlayerPieceCount = 0

		for nextX, nextY := x-1, y+1; nextX >= 0 && nextY < BoardHeight; nextX, nextY = nextX-1, nextY+1 {
			if gameState.board[nextY][nextX] != currentPlayerPiece {
				break
			}
			currentPlayerPieceCount++
		}

		for nextX, nextY := x+1, y-1; nextX < BoardWidth && nextY >= 0; nextX, nextY = nextX+1, nextY-1 {
			if gameState.board[nextY][nextX] != currentPlayerPiece {
				break
			}
			currentPlayerPieceCount++
		}

		if currentPlayerPieceCount >= 3 {
			currentPlayerWinOpportunities++
			continue
		}
	}

	if (heuristic.targetPlayer == Player1 && gameState.turn == Player1Turn) || (heuristic.targetPlayer == Player2 && gameState.turn == Player2Turn) {
		if currentPlayerWinOpportunities > 0 {
			// Heuristic player has opportunity to win this turn
			return 0.99
		}
	} else {
		// Opponent has opportunity to win this turn
		if heuristic.targetPlayer == Player2 {
			player1Viability += 10000 * currentPlayerWinOpportunities
		} else {
			player2Viability += 10000 * currentPlayerWinOpportunities
		}
	}

	// Check for horizontal viability
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

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

			switch gameState.board[y][x-3] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}
		}
	}

	// Check for vertical viability
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

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

			switch gameState.board[y-3][x] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}
		}
	}

	// Check for diagonally up viability
	for xIndex := 4 - BoardHeight; xIndex < BoardWidth-3; xIndex++ {
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

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

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

	// Check for diagonally down viability
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

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

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

	viability := float64(player1Viability-player2Viability) / float64(100+player1Viability+player2Viability)

	if heuristic.targetPlayer == Player1 {
		return viability
	} else {
		return -viability
	}
}
