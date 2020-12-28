package models

type Points []Point

func (points Points) Contains(other Point) bool {
	for _, p := range points {
		if p.Equal(other) {
			return true
		}
	}
	return false
}

