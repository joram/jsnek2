package main

import (
	"fmt"
	"github.com/BattlesnakeOfficial/rules"
	"github.com/joram/jsnek2/models"
	"math"
)

func rateState(current,previous models.MoveRequest) float64 {
	canReachOwnTrailRating := 10.0
	canReachOtherTrailRating := 1.0
	distanceToFoodRatingFactor := 100.0
	rating := 0.0

	if previous.IsSolid(current.You.Head()) {
		return -100
	}

	// can reach my tail
	_, err := current.Path(current.You.Head(), current.You.Tail())
	if err != nil {
		rating += canReachOwnTrailRating
	}

	// can reach other tails
	for _, snake := range current.OtherSnakes() {
		_, err := current.Path(current.You.Head(), snake.Tail())
		if err != nil {
			rating += canReachOtherTrailRating
		}
	}

	foundFood := false
	pathToFood := models.Points{}
	for _, food := range current.Board.Food {
		path, err := current.Path(current.You.Head(), models.Point(food))
		if err != nil {
			continue
		}
		if !foundFood || len(path) < len(pathToFood) {
			foundFood = true
			pathToFood = path
		}
	}
	if foundFood {
		rating += float64(100-len(pathToFood)) * distanceToFoodRatingFactor
	}

	return rating
}

func getNextMoveRequest(sr models.MoveRequest, direction string) models.MoveRequest {
	ruleset := rules.StandardRuleset{}
	moves := []rules.SnakeMove{
		{sr.You.ID, direction},
	}
	for _, snake := range sr.OtherSnakes() {
		moves = append(moves, rules.SnakeMove{ID:snake.ID, Move:rules.MoveUp})
	}

	nextBoardState, err := ruleset.CreateNextBoardState(&sr.Board, moves)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		panic("failed")
	}

	you := models.Snake{}
	for _, snake := range nextBoardState.Snakes {
		if snake.ID == sr.You.ID {
			you = models.Snake(snake)
			break
		}
	}

	return models.MoveRequest{
		Board: *nextBoardState,
		Turn:  sr.Turn + 1,
		You:   you,
	}
}

func move(sr models.MoveRequest) rules.SnakeMove {
	left := rateState(getNextMoveRequest(sr, rules.MoveLeft), sr)
	right := rateState(getNextMoveRequest(sr, rules.MoveRight), sr)
	up := rateState(getNextMoveRequest(sr, rules.MoveUp), sr)
	down := rateState(getNextMoveRequest(sr, rules.MoveDown), sr)

	max := math.Max(up, math.Max(down, math.Max(left, right)))

	direction := map[float64]string{
		up:    rules.MoveDown,
		down:  rules.MoveUp,
		left:  rules.MoveLeft,
		right: rules.MoveRight,
	}[max]

	fmt.Printf("up:%v, down:%v, left:%v, right:%v, max:%v dir:%s\n", up, down, left, right, max, direction)

	return rules.SnakeMove{
		Move: direction,
		ID:   sr.You.ID,
	}
}