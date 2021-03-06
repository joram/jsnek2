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

func (p *Point) Direction(to Point) string {
	xDelta := to.X - p.X
	if xDelta == 1 {
		return rules.MoveRight
	}
	if xDelta == -1 {
		return rules.MoveLeft
	}

	yDelta :=  to.Y - p.Y
	if yDelta == 1 {
		return rules.MoveUp
	}
	if yDelta == -1 {
		return rules.MoveDown
	}
	panic("non-adjcent square cant calculate direction")
}
