package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c *Coord) String() string {
	return fmt.Sprintf("%d_%d", c.X, c.Y)
}

type Snake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int     `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int     `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height           int                      `json:"height"`
	Width            int                      `json:"width"`
	Food             []Coord                  `json:"food"`
	Hazards          []Coord                  `json:"hazards"`
	Snakes           []Snake                  `json:"snakes"`
}

type Game struct {
	ID string `json:"id"`
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

func DecodeSnakeRequest(req *http.Request, decoded *SnakeRequest) error {
	err := json.NewDecoder(req.Body).Decode(&decoded)
	return err
}
