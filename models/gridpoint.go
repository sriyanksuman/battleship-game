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
		label := "A"
		if *w.AssignedPlayer == 2 {
			label = "B"
		}
		return fmt.Sprintf("%s-%s", label, w.AssignedObject.ID)
	}
	return "."
}
