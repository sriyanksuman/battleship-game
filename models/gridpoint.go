package models

import "fmt"

type Point struct {
	X int
	Y int
}

type Water struct {
	Point
	AssignedPlayer *int
	AssignedObject *Battleship
}

func NewWater(x, y int) *Water {
	return &Water{
		Point: Point{X: x, Y: y},
	}
}

func (w *Water) String() string {
	if w.AssignedObject != nil && w.AssignedPlayer != nil {
		return fmt.Sprintf("Player%d-BS%d(%d,%d)", *w.AssignedPlayer, w.AssignedObject.ID, w.X, w.Y)
	}
	if w.AssignedPlayer != nil {
		return fmt.Sprintf("Player%d(%d,%d)", *w.AssignedPlayer, w.X, w.Y)
	}
	return fmt.Sprintf("Water(%d,%d)", w.X, w.Y)
}
