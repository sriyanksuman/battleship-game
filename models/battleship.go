package models

type Battleship struct {
	ID         string
	Size       int
	Center     Point
	ShipPoints map[Point]bool
	IsActive   bool
}

func NewBattleship(id string, size int, center Point) *Battleship {
	bs := &Battleship{
		ID:         id,
		Size:       size,
		Center:     center,
		ShipPoints: make(map[Point]bool),
		IsActive:   true,
	}
	bs.populatePoints()
	return bs
}

func (bs *Battleship) populatePoints() {
	half := bs.Size / 2
	topLeftX := bs.Center.X - half
	topLeftY := bs.Center.Y - half
	for i := topLeftX; i < topLeftX+bs.Size; i++ {
		for j := topLeftY; j < topLeftY+bs.Size; j++ {
			bs.ShipPoints[Point{X: i, Y: j}] = true
		}
	}
}

func ValidateBattleship(gridSize, size int, center Point) bool {
	half := size / 2
	topLeftX := center.X - half
	topLeftY := center.Y - half
	x, y := topLeftX, topLeftY
	return x >= 0 && y >= 0 && x+size-1 < gridSize && y+size-1 < gridSize
}
