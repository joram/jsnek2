package models

func (mr *MoveRequest) OtherSnakes() []Snake {
	var otherSnakes []Snake
	for _, snake := range mr.Board.Snakes {
		if snake.ID != mr.You.ID {

			otherSnakes = append(otherSnakes, Snake{
				ID: snake.ID,
				Body: snake.Body,
			})
		}
	}
	return otherSnakes
}

func (mr *MoveRequest) IsEmpty(p Point) bool {
	// out of bounds
	if p.X < 0 || p.Y < 0 || p.X >= mr.Board.Width || p.Y >= mr.Board.Height {
		return false
	}

	// in snake body
	for _, snake := range mr.Board.Snakes {
		for _, b := range snake.Body {
			point := Point(b)
			if point.Equal(p){
				return false
			}
		}
	}
	return true
}

func (mr *MoveRequest) IsSolid(p Point) bool {
	return !mr.IsEmpty(p)
}
