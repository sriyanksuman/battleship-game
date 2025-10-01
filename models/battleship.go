package models

type Battleship struct {
	ID          int
	Size        int
	TopLeftEdge Point
	ShipPoints  map[Point]bool
	IsActive    bool
}

func NewBattleship(id, size int, topLeftEdge Point) *Battleship {
	bs := &Battleship{
		ID:          id,
		Size:        size,
		TopLeftEdge: topLeftEdge,
		ShipPoints:  make(map[Point]bool),
		IsActive:    true,
	}
	bs.populatePoints()
	return bs
}

func (bs *Battleship) populatePoints() {
	x, y := bs.TopLeftEdge.X, bs.TopLeftEdge.Y
	for i := x; i < x+bs.Size; i++ {
		for j := y; j < y+bs.Size; j++ {
			bs.ShipPoints[Point{X: i, Y: j}] = true
		}
	}
}

func ValidateBattleship(gridSize, size int, topLeftEdge Point) bool {
	x, y := topLeftEdge.X, topLeftEdge.Y
	return x >= 0 && y >= 0 && x+size-1 < gridSize && y+size-1 < gridSize
}
