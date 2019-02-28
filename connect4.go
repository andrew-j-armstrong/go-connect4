package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type Heuristic interface {
	Heuristic(game *Game) float64
}

func chooseHeuristic(player Connect4Player) Heuristic {
	for {
		fmt.Printf("Choose heurstic for player %d:\n", player)
		fmt.Printf("1: Simple\n")
		fmt.Printf("2: Viability\n")
		fmt.Printf("3: Viability Extended\n")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			return NewSimpleHeuristic(player)
		case 2:
			return NewViabilityHeuristic(player)
		case 3:
			return NewViabilityExtendedHeuristic(player)
		}

		fmt.Println("Invalid choice!")
	}
}

func chooseDifficulty() float64 {
	for {
		fmt.Printf("Choose difficulty (0-100): ")
		var choice float64
		fmt.Scan(&choice)

		if choice >= 0 && choice <= 100 {
			return choice
		}

		fmt.Println("Invalid choice!")
	}
}

func choosePlayer(game *Game, player Connect4Player) GameSimPlayer {
	for {
		fmt.Printf("Choose player %d:\n", player)
		fmt.Printf("1: Human\n")
		fmt.Printf("2: Randy\n")
		fmt.Printf("3: Huey\n")
		fmt.Printf("4: Max\n")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			return NewHumanPlayer(game)
		case 2:
			return NewRandomPlayer(game)
		case 3:
			return NewHeuristicPlayer(game, chooseHeuristic(player))
		case 4:
			difficulty := chooseDifficulty()

			var expectimaxPlayer *ExpectimaxPlayer
			switch player {
			case Player1:
				expectimaxPlayer = NewExpectimaxPlayer(game, ExpectimaxPlayer1, chooseHeuristic(player), difficulty, time.Duration(5)*time.Second)
			case Player2:
				expectimaxPlayer = NewExpectimaxPlayer(game, ExpectimaxPlayer2, chooseHeuristic(player), difficulty, time.Duration(5)*time.Second)
			}

			go expectimaxPlayer.Run()
			return expectimaxPlayer
		}

		fmt.Println("Invalid choice!")
	}
}

var playerDescriptionRegex = regexp.MustCompile(`(human)|(random)|(heuristic)/(simple|viability|viabilityextended)|(expectimax)/(simple|viability|viabilityextended)/(\d+)/(\d+)`)

func parseHeuristic(heuristicDescription string, player Connect4Player) Heuristic {
	switch heuristicDescription {
	case "simple":
		return NewSimpleHeuristic(player)
	case "viability":
		return NewViabilityHeuristic(player)
	case "viabilityextended":
		return NewViabilityExtendedHeuristic(player)
	default:
		log.Fatalf("invalid heuristic description %s", heuristicDescription)
		return nil
	}
}

func parseExpectimaxDifficulty(difficultyDescription string) float64 {
	difficulty, err := strconv.ParseFloat(difficultyDescription, 64)

	if err != nil {
		log.Fatal("error parsing expectimax difficulty ", err)
	}

	if difficulty < 0 || difficulty > 100 {
		log.Fatalf("invalid expectimax difficulty: %f is not between 0 and 100", difficulty)
	}

	return difficulty
}

func parsePlayer(game *Game, player Connect4Player, playerDescription string) GameSimPlayer {
	match := playerDescriptionRegex.FindStringSubmatch(playerDescription)
	if match == nil {
		log.Fatal("invalid player description", playerDescription)
	}

	if match[1] == "human" {
		fmt.Printf("Player %d: Human\n", player)
		return NewHumanPlayer(game)
	} else if match[2] == "random" {
		fmt.Printf("Player %d: Random\n", player)
		return NewRandomPlayer(game)
	} else if match[3] == "heuristic" {
		fmt.Printf("Player %d: Heuristic\n", player)
		return NewHeuristicPlayer(game, parseHeuristic(match[4], player))
	} else if match[5] == "expectimax" {
		fmt.Printf("Player %d: Expectimax\n", player)

		heuristic := parseHeuristic(match[6], player)
		difficulty := parseExpectimaxDifficulty(match[7])
		maxDurationMilliseconds, _ := strconv.Atoi(match[8])
		maxDuration := time.Duration(maxDurationMilliseconds) * time.Millisecond

		var expectimaxPlayer *ExpectimaxPlayer
		switch player {
		case Player1:
			expectimaxPlayer = NewExpectimaxPlayer(game, ExpectimaxPlayer1, heuristic, difficulty, maxDuration)
		case Player2:
			expectimaxPlayer = NewExpectimaxPlayer(game, ExpectimaxPlayer2, heuristic, difficulty, maxDuration)
		}

		go expectimaxPlayer.Run()
		return expectimaxPlayer
	} else {
		log.Fatalf("unknown player type %s", playerDescription)
		return nil
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	runGamePtr := flag.Bool("run", false, "Run a game")
	loadGamePtr := flag.String("load", "", "Load a game from <filename>")
	saveGamePtr := flag.String("save", "", "Save the game after the move to <filename>")
	player1Ptr := flag.String("player1", "", "Player 1 player description")
	player2Ptr := flag.String("player2", "", "Player 2 player description")
	makeMovePtr := flag.String("makeMove", "", "Make this move")
	nextPlayerPtr := flag.String("nextPlayer", "", "Next player description")

	flag.Parse()

	var game *Game
	var player1 GameSimPlayer
	var player2 GameSimPlayer

	if (*loadGamePtr) == "" {
		game = NewGame()
	} else {
		// Load from file
		var err error
		game, err = LoadGame(*loadGamePtr)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Loaded game:")
		game.Print()
	}

	if game.IsGameOver() {
		game.Print()
		return
	}

	if *makeMovePtr != "" {
		move := ParseMove(*makeMovePtr)
		if !game.IsValidMove(move) {
			log.Fatalf("invalid move: %s", move.String())
		}

		err := game.MakeMove(move)
		if err != nil {
			log.Fatal(err)
		}

		game.Print()
	}

	if game.IsGameOver() {
		game.Print()
		return
	}

	if *runGamePtr {
		if *player1Ptr == "" {
			player1 = choosePlayer(game, 1)
		} else {
			player1 = parsePlayer(game, 1, *player1Ptr)
		}

		if *player2Ptr == "" {
			player2 = choosePlayer(game, 2)
		} else {
			player2 = parsePlayer(game, 2, *player2Ptr)
		}

		gameSim := NewGameSim(game, player1, player2)
		gameSim.Run()
		gameSim.game.Print()
	} else if *nextPlayerPtr != "" || (game.turn == Player1Turn && *player1Ptr != "") || (game.turn == Player2Turn && *player2Ptr != "") {
		// Perform one move
		var player GameSimPlayer
		switch game.turn {
		case Player1Turn:
			if *nextPlayerPtr != "" {
				player = parsePlayer(game, 1, *nextPlayerPtr)
			} else if *player1Ptr == "" {
				player = choosePlayer(game, 1)
			} else {
				player = parsePlayer(game, 1, *player1Ptr)
			}
		case Player2Turn:
			if *nextPlayerPtr != "" {
				player = parsePlayer(game, 2, *nextPlayerPtr)
			} else if *player2Ptr == "" {
				player = choosePlayer(game, 2)
			} else {
				player = parsePlayer(game, 2, *player2Ptr)
			}
		}

		switch player.(type) {
		case *ExpectimaxPlayer:
			time.Sleep(time.Duration(5) * time.Second)
		}

		game.MakeMove(player.GetNextMove())
		game.Print()

		if *saveGamePtr != "" {
			game.Save(*saveGamePtr)
		}
	} else {
		game.Print()
	}
}
