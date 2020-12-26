package main

import (
	"encoding/json"
	"github.com/BattlesnakeOfficial/rules"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)


func Start(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	respond(res, StartResponse{
		APIVersion: "1",
		Author: "John Oram",
		Color: "#75CEDD",
		Head: "silly",
		Tail: "fat-rattle",
	})
}

type SnakeRequest struct {
	Turn  int   `json:"turn"`
	Board rules.BoardState `json:"board"`
	You   rules.Snake `json:"you"`
}

func Move(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sr := SnakeRequest{}
	err := json.NewDecoder(req.Body).Decode(&sr)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}

	response := MoveResponse{Move: move(sr).Move}
	respond(res, response)
}

func End(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sr := SnakeRequest{}
	err := json.NewDecoder(req.Body).Decode(&sr)
	if err != nil {
		log.Printf("Bad end request: %v", err)
	}
}

func Ping(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
}
