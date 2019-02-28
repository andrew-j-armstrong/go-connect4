package main

import (
	"fmt"
	"math"
	"testing"

	"github.com/carbon-12/go-extensions"
)

func TestCalculateChildLikelihood(t *testing.T) {
	player := ExpectimaxPlayer{}

	t.Run("test calculateChildLikelihood()", func(t *testing.T) {
		testCases := []struct {
			valueDifference      float64
			expectimaxDifficulty float64
			minSpread            float64
			minMove1Likelihood   float64
			maxMove1Likelihood   float64
		}{
			{0.001, 100.0, 0.0, 0.99, 1.0},
			{1.0, 0.0, 0.0, 0.50, 0.51},
			{0.1, 75.0, 0.0, 0.95, 0.99},
			{0.1, 50.0, 0.0, 0.80, 0.90},
			{0.1, 25.0, 0.0, 0.55, 0.70},
			{0.2, 75.0, 0.0, 0.99, 0.9999},
			{0.2, 50.0, 0.0, 0.90, 0.98},
			{0.2, 25.0, 0.0, 0.6, 0.80},
			{0.3, 75.0, 0.0, 0.999, 1.0},
			{0.3, 50.0, 0.0, 0.98, 0.99},
			{0.3, 25.0, 0.0, 0.6, 0.80},
			{0.4, 75.0, 0.0, 0.999, 1.0},
			{0.4, 50.0, 0.0, 0.99, 0.999},
			{0.4, 25.0, 0.0, 0.6, 0.80},
		}

		for _, testCase := range testCases {
			childLikelihood := extensions.ValueMap{1: 0.0, 2: 0.0}

			player.calculateChildLikelihoodMap(func(move interface{}) float64 {
				if move == 1 {
					return 1.0
				} else {
					return 1.0 - testCase.valueDifference
				}
			}, &childLikelihood, testCase.expectimaxDifficulty, testCase.minSpread)

			if childLikelihood[1] < testCase.minMove1Likelihood {
				t.Error(fmt.Sprintf("move likelihood too low (%f%%) ", 100.0*childLikelihood[1]), testCase)
			} else if childLikelihood[1] > testCase.maxMove1Likelihood {
				t.Error(fmt.Sprintf("move likelihood too high (%f%%) ", 100.0*childLikelihood[1]), testCase)
			} else {
				childLikelihood1 := childLikelihood[1]

				player.calculateChildLikelihoodMap(func(move interface{}) float64 {
					if move == 1 {
						return 0.0
					} else {
						return 0.0 - testCase.valueDifference
					}
				}, &childLikelihood, testCase.expectimaxDifficulty, testCase.minSpread)

				if math.Abs(childLikelihood1-childLikelihood[1]) > 0.00001 {
					t.Error(fmt.Sprintf("move likelihood differs at different ranges (%f%% for 1.0, %f%% for 0.0) ", 100.0*childLikelihood1, 100.0*childLikelihood[1]), testCase)
				}
			}
		}
	})
}
