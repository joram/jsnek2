package models

import (
	"github.com/BattlesnakeOfficial/rules"
)
type Snake rules.Snake

func (snake *Snake) Head() Point {
	return Point(snake.Body[0])
}

func (snake *Snake) Tail() Point {
	return Point(snake.Body[len(snake.Body)-1])
}
