package models

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
	if w.AssignedObject != nil {
		return w.AssignedObject.ID
	}
	return "."
}
