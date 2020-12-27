package models

import (
	"errors"
	"github.com/BattlesnakeOfficial/rules"
)


type MoveRequest struct {
	Turn  int        `json:"turn"`
	Board rules.BoardState `json:"board"`
	You   Snake      `json:"you"`
}

type WeightedPoint struct {
	weight int
	point Point
}

func (mr *MoveRequest) Path(from, to Point) ([]Point, error) {
	visited := Points{}
	toVisit := []WeightedPoint{
		{0, to},
	}
	weight := map[Point]int {}
	for {
		visiting, toVisit := toVisit[0], toVisit[1:]
		weight[visiting.point] = visiting.weight

		for _, adjacentPoint := range visiting.point.Adjacent() {
			if visited.contains(adjacentPoint) {
				continue
			}
			if (mr.IsEmpty(adjacentPoint) || adjacentPoint.Equal(from)) && !visited.contains(adjacentPoint) {
				toVisit = append(toVisit, WeightedPoint{visiting.weight, adjacentPoint})
				visited = append(visited, adjacentPoint)
			}
		}

		if len(toVisit) == 0 || visiting.point == from {
			break
		}
	}

	if !visited.contains(from) {
		return Points{}, errors.New("no path found")
	}

	path := Points{from}
	currentWeight := weight[from]
	currentPoint := from
	for {
		for _, adjacentPoint := range currentPoint.Adjacent() {
			if adjacentPoint.Equal(to) {
				path = append(path, to)
				return path, nil
			}
			if weight[adjacentPoint] == currentWeight - 1 {
				path = append(path, adjacentPoint)
				currentWeight -= 1
				currentPoint = adjacentPoint
				break
			}
		}
		if currentPoint.Equal(to) {
			path = append(path, to)
			return path, nil
		}
	}
}


type StartResponse struct {
	APIVersion string `json:"apiversion,omitempty"`
	Author string `json:"author,omitempty"`
	Color string `json:"color,omitempty"`
	Head string `json:"head,omitempty"`
	Tail string `json:"tail,omitempty"`
}


type MoveResponse struct {
	Move  string `json:"move"`
	Taunt string `json:"taunt"`
}
