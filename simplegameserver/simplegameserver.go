package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zerok/go-dojo-pairs/game"
)

type gameState struct {
	SolvedCards   []string
	Turned        int
	TurnedCard    string
	CurrentPlayer byte
	GameOver      bool
	Scores        [2]int
}

func main() {
	var g game.PairsGame
	var numberOfPairs int
	var listenAddr string

	flag.StringVar(&listenAddr, "listen", "localhost:8888", "Address the webserver should listen on")
	flag.IntVar(&numberOfPairs, "pairs", 10, "Number of pairs in the deck")

	flag.Parse()

	if numberOfPairs < 2 {
		log.Fatalln("Please require at least 2 pairs.")
	}

	g, _ = game.NewGame(numberOfPairs)
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("./static")))
	router.HandleFunc("/api/status/", func(w http.ResponseWriter, r *http.Request) {
		turnedCard, turnedIndex := g.TurnedCard()
		s := gameState{
			GameOver:      false,
			Scores:        g.Scores(),
			CurrentPlayer: g.CurrentPlayer(),
			Turned:        turnedIndex,
			TurnedCard:    turnedCard,
			SolvedCards:   g.SolvedCards(),
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(s)
	})
	router.HandleFunc("/api/pick/{card:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		card, ok := vars["card"]
		if !ok {
			http.NotFound(w, r)
			return
		}
		cardIdx, err := strconv.Atoi(card)
		if err != nil {
			http.Error(w, "Not a number", 400)
			return
		}
		s, err := g.TurnCard(cardIdx, g.CurrentPlayer())
		if err != nil {
			log.Printf("Turn failed: %s", err.Error())
			http.Error(w, "Turn failed", 400)
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(s)
	})
	log.Printf("Starting server on http://%s", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
