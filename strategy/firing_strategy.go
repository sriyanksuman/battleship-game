package strategy

import (
	"battleship-game/models"
	"math/rand"
	"time"
)

type FiringStrategy interface {
	Fire(currentPlayer *models.Player, battlefield *models.Battlefield, previousPoint models.Point) *models.Water
}

type RandomFiringStrategy struct {
	SampleSet []models.Point
	rng       *rand.Rand
}

func NewRandomFiringStrategy(sampleSet []models.Point) *RandomFiringStrategy {
	source := rand.NewSource(time.Now().UnixNano())
	return &RandomFiringStrategy{
		SampleSet: sampleSet,
		rng:       rand.New(source),
	}
}

func (rfs *RandomFiringStrategy) Fire(currentPlayer *models.Player, battlefield *models.Battlefield, previousPoint models.Point) *models.Water {
	var point *models.Water
	var i, j int

	for {
		index := rfs.rng.Intn(len(rfs.SampleSet))
		i = rfs.SampleSet[index].X
		j = rfs.SampleSet[index].Y
		point = battlefield.Grid[i][j]

		if (point.AssignedPlayer == nil || *point.AssignedPlayer != currentPlayer.ID) &&
			!(i == previousPoint.X && j == previousPoint.Y) {
			break
		}
	}

	return point
}
