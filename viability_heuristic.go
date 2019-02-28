package main

type ViabilityHeuristic struct {
	targetPlayer Connect4Player
}

func NewViabilityHeuristic(targetPlayer Connect4Player) *ViabilityHeuristic {
	return &ViabilityHeuristic{targetPlayer}
}

func (heuristic *ViabilityHeuristic) increaseViabilityScores(player1PieceCount int, player2PieceCount int, player1Viability *int, player2Viability *int) {
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

func (heuristic *ViabilityHeuristic) Heuristic(game *Game) float64 {
	if game.turn == Draw {
		return 0.0
	} else if game.turn == Player1Won {
		if heuristic.targetPlayer == Player1 {
			return 1.0
		} else {
			return -1.0
		}
	} else if game.turn == Player2Won {
		if heuristic.targetPlayer == Player1 {
			return -1.0
		} else {
			return 1.0
		}
	}

	var player1Viability int
	var player2Viability int

	// Check for horizontal viability
	for y := 0; y < BoardHeight; y++ {
		player1PieceCount := 0
		player2PieceCount := 0

		x := 0
		for ; x < 3; x++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}

		for ; x < BoardWidth; x++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

			switch game.board[y][x-3] {
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
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}

		for ; y < BoardHeight; y++ {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

			switch game.board[y-3][x] {
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
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			x++
			y++
		}

		for x < BoardWidth && y < BoardHeight {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

			switch game.board[y-3][x-3] {
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
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			x++
			y--
		}

		for x < BoardWidth && y >= 0 {
			switch game.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}

			heuristic.increaseViabilityScores(player1PieceCount, player2PieceCount, &player1Viability, &player2Viability)

			switch game.board[y+3][x-3] {
			case Player1Piece:
				player1PieceCount--
			case Player2Piece:
				player2PieceCount--
			}

			x++
			y--
		}
	}

	var viability float64 = float64(100+player1Viability-player2Viability) / float64(200+player1Viability+player2Viability)

	if heuristic.targetPlayer == Player1 {
		return viability
	} else {
		return -viability
	}
}
