package models

import "github.com/BattlesnakeOfficial/rules"

type Point rules.Point

func (p *Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p *Point) Adjacent() Points {
	return Points{
		{p.X+1, p.Y},
		{p.X-1, p.Y},
		{p.X, p.Y+1},
		{p.X, p.Y-1},
	}
}
