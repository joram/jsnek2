package main

import (
	"encoding/json"
	"fmt"
	"github.com/joram/jsnek2/models"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)


func Start(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	respond(res, models.StartResponse{
		APIVersion: "1",
		Author: "John Oram",
		Color: "#ff6666",
		Head: "silly",
		Tail: "fat-rattle",
	})
}

func Move(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sr := models.MoveRequest{
		WeightedMaps: map[string]models.WeightedMap{},
		WeightedMapsAchievedGoal: map[string]bool{},
	}
	err := json.NewDecoder(req.Body).Decode(&sr)


	if err != nil {
		fmt.Printf("Bad move request: %v\n", err)
	}

	response := models.MoveResponse{Move: move(sr).Move}
	respond(res, response)
}

func End(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sr := models.MoveRequest{}
	err := json.NewDecoder(req.Body).Decode(&sr)
	if err != nil {
		log.Printf("Bad end request: %v", err)
	}
}

func Ping(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
}
