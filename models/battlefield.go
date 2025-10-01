package models

import (
	"fmt"
	"strings"
)

type Battlefield struct {
	N    int
	Grid [][]*Water
}

func NewBattlefield(n int) *Battlefield {
	grid := make([][]*Water, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]*Water, n)
		for j := 0; j < n; j++ {
			grid[i][j] = NewWater(i, j)
		}
	}
	return &Battlefield{
		N:    n,
		Grid: grid,
	}
}

func (bf *Battlefield) PrintBattlefield() {
	width := 20
	borderParts := make([]string, bf.N)
	for i := 0; i < bf.N; i++ {
		borderParts[i] = strings.Repeat("-", width)
	}
	border := "+" + strings.Join(borderParts, "+") + "+"

	fmt.Println(border)
	for i := 0; i < bf.N; i++ {
		row := bf.Grid[i]
		cellStrs := make([]string, bf.N)
		for j := 0; j < bf.N; j++ {
			cellStr := row[j].String()
			if len(cellStr) < width {
				cellStr = cellStr + strings.Repeat(" ", width-len(cellStr))
			}
			cellStrs[j] = cellStr
		}
		fmt.Println("|" + strings.Join(cellStrs, "|") + "|")
	}
	fmt.Println(border)
}
