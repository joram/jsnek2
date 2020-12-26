package main

import (
	"github.com/BattlesnakeOfficial/rules"
)

func move(sr SnakeRequest) rules.SnakeMove {
	return rules.SnakeMove{
		Move: rules.MoveLeft,
		ID: sr.You.ID,
	}
}
