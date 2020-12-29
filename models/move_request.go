package models

import (
	"errors"
	"fmt"
	"github.com/BattlesnakeOfficial/rules"
	"reflect"
)

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
		for _, b := range snake.Body[:len(snake.Body)-1] {
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


type WeightedMap map[Point]int
type WeightedMapInput struct {
	from Points
	goals Points
}

func (wmi *WeightedMapInput) HashKey() string {
	s := ""
	for _, p := range wmi.from {
		s = fmt.Sprintf("%s%v", s, p)
	}
	s = fmt.Sprintf("%s:", s)
	for _, p := range wmi.goals {
		s = fmt.Sprintf("%s%v", s, p)
	}
	return s
}

type MoveRequest struct {
	Turn  int              `json:"turn"`
	Board rules.BoardState `json:"board"`
	You   Snake            `json:"you"`

	WeightedMaps             map[string]WeightedMap
	WeightedMapsAchievedGoal map[string]bool
}

type WeightedPoint struct {
	weight int
	point Point
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func (mr *MoveRequest) WeightedMap(from, goals Points) (WeightedMap, bool) {

	// get cached
	input := WeightedMapInput{from, goals}
	key := input.HashKey()

	value, exists := mr.WeightedMaps[key]
	if exists {
		return value, mr.WeightedMapsAchievedGoal[key]
	}

	visited := 0
	gonnaVisit := map[Point]bool{}
	heat := map[Point]int{}
	var toVisit []WeightedPoint
	for _, p := range from {
		tv := WeightedPoint{0,p}
		if !gonnaVisit[p]{
			toVisit = append(toVisit, tv)
			gonnaVisit[p] = true
			heat[p] = 0
		}
	}
	for {
		// done
		if len(toVisit) == 0 {
			break
		}

		// pop
		tv := toVisit[0]
		toVisit = toVisit[1:]

		// visit
		heat[tv.point] = tv.weight
		if goals.Contains(tv.point) {
			mr.WeightedMapsAchievedGoal[key] = true
			mr.WeightedMaps[key] = heat
			return heat, true
		}

		// visit next
		for _, next := range tv.point.Adjacent() {
			if !gonnaVisit[next] && mr.IsEmpty(next) {
				toVisit = append(toVisit, WeightedPoint{tv.weight + 1, next})
				gonnaVisit[next] = true
			}
		}
		visited += 1

	}
	mr.WeightedMapsAchievedGoal[key] = false
	mr.WeightedMaps[key] = heat

	return heat, false
}

func (mr *MoveRequest) CanAccess(from Point) int {
	weights, _ := mr.WeightedMap(Points{from}, Points{})
	return len(weights)
}

func (mr *MoveRequest) Path(from, to Point) ([]Point, error) {
	weights, reachedGoal := mr.WeightedMap(Points{from}, Points{to})
	if !reachedGoal {
		return Points{}, errors.New("no path found")
	}

	path := Points{to}
	currentPoint := to
	for !currentPoint.Equal(from) {
		for _, adjacentPoint := range currentPoint.Adjacent() {
			weight, exists := weights[adjacentPoint]
			if !exists {
				continue
			}
			if weight == weights[currentPoint] - 1 {
				path = append(path, adjacentPoint)
				currentPoint = adjacentPoint
				break
			}
		}
	}
	reverse(path)
	return path, nil
}

func (mr *MoveRequest) EmptyPoints(points Points) Points {
	var empties Points
	for _, p := range points {
		if mr.IsEmpty(p) {
			empties = append(empties, p)
		}
	}
	return empties
}

func (mr *MoveRequest) IsDead(id string) bool {
	for _, snake := range mr.Board.Snakes {
		if snake.ID == id {
			return false
		}
	}
	return true
}

func (mr *MoveRequest) IsEdge(p Point) bool {
	if p.X == 0 || p.Y == 0 {
		return true
	}
	if p.X == mr.Board.Width-2 || p.Y == mr.Board.Height-2 {
		return true
	}
	return false
}
