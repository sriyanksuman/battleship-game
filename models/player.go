package models

import "fmt"

type Player struct {
	ID                int
	Name              string
	Battleships       []*Battleship
	ActiveBattleships map[int]bool
}

func NewPlayer(id int, name string) *Player {
	return &Player{
		ID:                id,
		Name:              name,
		Battleships:       make([]*Battleship, 0),
		ActiveBattleships: make(map[int]bool),
	}
}

func (p *Player) AddBattleship(bs *Battleship) {
	p.Battleships = append(p.Battleships, bs)
	p.ActiveBattleships[bs.ID] = true
}

func (p *Player) String() string {
	activeShips := make([]int, 0)
	for shipID := range p.ActiveBattleships {
		activeShips = append(activeShips, shipID)
	}
	return fmt.Sprintf("Player%d, activeBattleships: %v", p.ID, activeShips)
}
