package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"
	"time"

	"github.com/carbon-12/go-expectimax"
	"github.com/carbon-12/go-extensions"
)

type ExpectimaxPlayerID int

const (
	ExpectimaxPlayer1 ExpectimaxPlayerID = iota
	ExpectimaxPlayer2
)

type ExpectimaxHeuristic interface {
	Heuristic(game *Game) float64
}

func getExpectimaxHeuristic(heuristic ExpectimaxHeuristic) expectimax.ExpectimaxHeuristic {
	return func(expectimaxGame expectimax.Game) float64 {
		game, ok := expectimaxGame.(*Game)

		if !ok {
			log.Fatal(errors.New("Connect4 ExpectimaxHeuristic called for invalid game!"))
			return 0.0
		}

		return heuristic.Heuristic(game)
	}
}

type Expectimax interface {
	RunExpectimax()
	IsCurrentlySearching() bool
	GetNextMoveValues() *extensions.ValueMap
}

type ExpectimaxPlayer struct {
	game           *Game
	player         ExpectimaxPlayerID
	expectimaxBase Expectimax
	difficulty     float64
	lastChoiceTime time.Time
	maxSearchTime  time.Duration
}

func (player *ExpectimaxPlayer) Run() {
	player.lastChoiceTime = time.Now()
	player.expectimaxBase.RunExpectimax()
}

func (player *ExpectimaxPlayer) IsReadyToMakeMove() bool {
	if !player.expectimaxBase.IsCurrentlySearching() {
		return true
	}

	nextMoves := player.expectimaxBase.GetNextMoveValues()

	// Check the moves to determine whether there's enough difference or we should wait longer
	if len(*nextMoves) <= 1 {
		return true
	}

	moveValues := make([]float64, 0, len(*nextMoves))
	for _, value := range *nextMoves {
		moveValues = append(moveValues, value)
	}

	sort.Float64s(moveValues)

	return moveValues[len(moveValues)-1] >= 1.0 || moveValues[len(moveValues)-1] <= -1.0 || moveValues[len(moveValues)-1]-moveValues[len(moveValues)-2] >= 0.5*player.difficulty/100.0
}

func (player *ExpectimaxPlayer) buildSelectionWheel(moves *extensions.InterfaceSlice, getMoveValue func(interface{}) float64, playerDifficulty float64) *extensions.ValueMap {
	powerBase := math.Pow(100.0/(100.0-playerDifficulty), playerDifficulty/25.0)
	selectionWheel := extensions.ValueMap{}
	for _, move := range *moves {
		wheelValue := math.Pow(powerBase, 10.0*getMoveValue(move))
		selectionWheel[move] = wheelValue
	}

	return &selectionWheel
}

func (player *ExpectimaxPlayer) GetNextMove() Move {
	player.lastChoiceTime = time.Now()

	for time.Since(player.lastChoiceTime) < player.maxSearchTime && !player.IsReadyToMakeMove() {
		time.Sleep(time.Duration(50) * time.Millisecond)
	}

	nextMoves := player.expectimaxBase.GetNextMoveValues()

	if len(*nextMoves) == 0 {
		log.Fatal("No next moves!")
	}

	var nextMove Move
	if player.difficulty == 100.0 {
		nextMove = nextMoves.GetBestKey().(Move)
	} else {
		selectionWheel := player.buildSelectionWheel(nextMoves.GetKeys(), func(move interface{}) float64 { return (*nextMoves)[move] }, player.difficulty)
		nextMove = selectionWheel.SelectFromWheel().(Move)
	}

	player.lastChoiceTime = time.Now()

	return nextMove
}

func (player *ExpectimaxPlayer) calculateChildLikelihoodMap(getChildValue func(interface{}) float64, childLikelihood *extensions.ValueMap, playerDifficulty float64, minSpread float64) {
	selectionWheel := player.buildSelectionWheel(childLikelihood.GetKeys(), getChildValue, playerDifficulty)
	totalWheelValue := selectionWheel.GetTotalValue()

	for move, wheelValue := range *selectionWheel {
		(*childLikelihood)[move] = (minSpread / float64(len(*selectionWheel))) + ((1.0 - minSpread) * wheelValue / totalWheelValue)
	}
}

func (player *ExpectimaxPlayer) calculateChildLikelihood(getGame func() expectimax.Game, getChildValue func(interface{}) float64, childLikelihood *extensions.ValueMap) {
	genericGame := getGame()
	if genericGame == nil {
		return
	}

	game, ok := genericGame.(*Game)

	if !ok {
		log.Fatal(fmt.Sprintf("calculateChildLikelihood received invalid game type: %s", reflect.TypeOf(genericGame)))
	}

	if game.turn == Player1Turn && player.player == ExpectimaxPlayer2 || game.turn == Player2Turn && player.player == ExpectimaxPlayer1 {
		player.calculateChildLikelihoodMap(func(move interface{}) float64 { return -getChildValue(move) }, childLikelihood, 99, 0.02)
	} else {
		if player.difficulty == 100.0 {
			bestMove := childLikelihood.GetKeys().GetBestEntry(getChildValue)

			for move := range *childLikelihood {
				if move == bestMove {
					(*childLikelihood)[move] = 1.0
				} else {
					(*childLikelihood)[move] = 0.0
				}
			}
		} else {
			player.calculateChildLikelihoodMap(getChildValue, childLikelihood, player.difficulty, 0.0)
		}
	}
}

func NewExpectimaxPlayer(game *Game, playerID ExpectimaxPlayerID, heuristic ExpectimaxHeuristic, difficulty float64, maxSearchTime time.Duration) *ExpectimaxPlayer {
	player := &ExpectimaxPlayer{game, playerID, nil, difficulty, time.Time{}, maxSearchTime}
	player.expectimaxBase = expectimax.NewExpectimax(game, getExpectimaxHeuristic(heuristic), player.calculateChildLikelihood, 10000)
	return player
}
