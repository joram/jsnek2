package main

import (
	"context"
	"fmt"
	"github.com/BattlesnakeOfficial/rules"
	"github.com/joram/jsnek2/models"
	"math"
	"math/rand"
	"time"
)

func rateState(current, previous models.MoveRequest) float64 {
	canReachOwnTrailRating := 10.0
	canReachOtherTrailRating := 1.0
	distanceToFoodRatingFactor := 10.0
	potentialHeadToHeadLossRating := -100.0
	potentialHeadToHeadWinRating := 100.0
	isSolidRating := -1000.0
	isDeadRating := -1000.0
	isLastRemaining := 1000.0
	isEdgeRating := -100.0

	rating := 0.0

	// is last remaining
	if len(current.Board.Snakes) == 1 && current.Board.Snakes[0].ID == current.You.ID {
		rating += isLastRemaining
	}

	rating += float64(current.You.Health)

	if current.IsEdge(current.You.Head()) {
		rating += isEdgeRating
	}

	// is dead
	isAlive := false
	for _, snake := range current.Board.Snakes {
		if snake.ID == current.You.ID {
			isAlive = true
		}
	}
	if !isAlive {
		return isDeadRating
	}

	// solid wall
	if previous.IsSolid(current.You.Head()) {
		return isSolidRating
	}

	// can reach my tail
	_, err := current.Path(current.You.Head(), current.You.Tail())
	if err != nil {
		rating += canReachOwnTrailRating
	}

	// can head-to-head
	for _, snake := range previous.OtherSnakes() {
		head := snake.Head()
		if head.Adjacent().Contains(current.You.Head()) {
			if len(snake.Body) >= len(current.You.Body) {
				rating += potentialHeadToHeadLossRating
			} else {
				rating += potentialHeadToHeadWinRating
			}
		}
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

func otherSnakeMove(sr models.MoveRequest, snake models.Snake) string {
	head := snake.Head()

	empties := sr.EmptyPoints(head.Adjacent())
	if len(empties) == 1 {
		return head.Direction(empties[0])
	}

	//path, err := sr.Path(head, sr.You.Head())
	//if err == nil {
	//	return head.Direction(path[1])
	//}

	if len(empties) > 1 {
		i := rand.Intn(len(empties))
		return head.Direction(empties[i])
	}

	choices := []string{
		rules.MoveLeft,
		rules.MoveRight,
		rules.MoveUp,
		rules.MoveDown,
	}
	return choices[rand.Intn(len(choices))]

}

func getNextMoveRequest(mr models.MoveRequest, direction string) models.MoveRequest {
	ruleset := rules.StandardRuleset{}

	// moves
	moves := []rules.SnakeMove{
		{mr.You.ID, direction},
	}
	for _, snake := range mr.OtherSnakes() {
		moves = append(moves, rules.SnakeMove{
			ID:   snake.ID,
			Move: otherSnakeMove(mr, snake),
		})
	}

	nextBoardState, err := ruleset.CreateNextBoardState(&mr.Board, moves)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		panic("failed")
	}

	you := models.Snake{}
	for _, snake := range nextBoardState.Snakes {
		if snake.ID == mr.You.ID {
			you = models.Snake(snake)
			break
		}
	}

	return models.MoveRequest{
		Board:                    *nextBoardState,
		Turn:                     mr.Turn + 1,
		You:                      you,
		WeightedMaps:             map[string]models.WeightedMap{},
		WeightedMapsAchievedGoal: map[string]bool{},
	}
}

func move(sr models.MoveRequest) rules.SnakeMove {
	ctx, cancel := context.WithCancel(context.Background())

	leftMoveRequest := getNextMoveRequest(sr, rules.MoveLeft)
	rightMoveRequest := getNextMoveRequest(sr, rules.MoveRight)
	upMoveRequest := getNextMoveRequest(sr, rules.MoveUp)
	downMoveRequest := getNextMoveRequest(sr, rules.MoveDown)

	leftChan := make(chan float64)
	rightChan := make(chan float64)
	upChan := make(chan float64)
	downChan := make(chan float64)

	go findBestWorstCaseScenarios(rules.MoveLeft, leftMoveRequest, sr, leftChan, ctx)
	go findBestWorstCaseScenarios(rules.MoveRight, rightMoveRequest, sr, rightChan, ctx)
	go findBestWorstCaseScenarios(rules.MoveUp, upMoveRequest, sr, upChan, ctx)
	go findBestWorstCaseScenarios(rules.MoveDown, downMoveRequest, sr, downChan, ctx)

	time.Sleep(200 * time.Millisecond)
	cancel()

	worstLeft,bestLeft := <-leftChan,<-leftChan
	worstright, bestRight := <-rightChan,<-rightChan
	worstUp, bestUp := <-upChan, <-upChan
	worstDown, bestDown := <-downChan,<-downChan

	bestWorst := math.Max(worstLeft, math.Max(worstright, math.Max(worstUp, worstDown)))
	bestBest := math.Max(bestLeft, math.Max(bestRight, math.Max(bestUp, bestDown)))

	direction := map[float64]string{
		bestUp:    rules.MoveDown,
		bestDown:  rules.MoveUp,
		bestLeft:  rules.MoveLeft,
		bestRight: rules.MoveRight,
	}[bestBest]

	fmt.Printf("up:%v, down:%v, left:%v, right:%v, bestWorst:%v bestBest:%v, dir:%s\n", bestUp, bestDown, bestLeft, bestRight, bestWorst, bestBest, direction)

	return rules.SnakeMove{
		Move: direction,
		ID:   sr.You.ID,
	}
}

func findBestWorstCaseScenarios(dir string, current, previous models.MoveRequest, c chan float64, ctx context.Context) {
	i := 0
	rating := rateState(current, previous)
	worstRating := rating
	bestRating := rating
	path := models.Points{
		previous.You.Head(),
	}
	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			break
		default:
			if current.IsDead(current.You.ID) {
				done = true
				break
			}

			head := current.You.Head()
			if path.Contains(head) {
				done = true
				break
			}

			i += 1
			leftMoveRequest := getNextMoveRequest(current, rules.MoveLeft)
			rightMoveRequest := getNextMoveRequest(current, rules.MoveRight)
			upMoveRequest := getNextMoveRequest(current, rules.MoveUp)
			downMoveRequest := getNextMoveRequest(current, rules.MoveDown)
			left := rateState(leftMoveRequest, current)
			right := rateState(rightMoveRequest, current)
			up := rateState(upMoveRequest, current)
			down := rateState(downMoveRequest, current)
			max := math.Max(up, math.Max(down, math.Max(left, right)))
			worstRating = math.Min(worstRating, max)
			bestRating = math.Max(bestRating, max)

			direction := map[float64]string{
				up:    rules.MoveDown,
				down:  rules.MoveUp,
				left:  rules.MoveLeft,
				right: rules.MoveRight,
			}[max]

			previous = current
			current = getNextMoveRequest(current, direction)

			path = append(path, head)
			if len(path) > 50 {
				done = true
				break
			}

			time.Sleep(1 * time.Nanosecond)
		}
	}
	c <- worstRating
	c <- bestRating
	fmt.Printf("%s, i:%d, worst:%f, best:%f\n", dir, i, worstRating, bestRating)
}
