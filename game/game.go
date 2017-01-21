package game

import (
	"fmt"
	"math/rand"

	"github.com/docker/docker/pkg/namesgenerator"
)

const maxNumberOfPairs = 10

// PairsGame represents an interface that should be implemented during the
// dojo in order for the game to be used by the simplegameserver application.
type PairsGame interface {
	TurnedCard() (string, int)
	Scores() [2]int
	SolvedCards() []string
	CurrentPlayer() byte
	TurnCard(i int, player byte) (TurnResult, error)
}

// TurnResult is the outcome of TurnCard operation representing the new state
// afterwards.
type TurnResult struct {
	Player         byte
	NextTurnPlayer byte
	TurnComplete   bool
	MatchFound     bool
	NewScore       int
	Card           string
	GameOver       bool
	Winner         byte
}

// Game abstracts a game of pairs. After initializing it, you can take turns by
// turning cards.
type Game struct {
	cards           []string
	turnedCardIndex int
	currentPlayer   byte
	scores          [2]int
	pairsLeft       int
	solvedCards     []string
}

func (g *Game) TurnedCard() (string, int) {
	c, _ := g.getCard(g.turnedCardIndex)
	return c, g.turnedCardIndex
}

func (g *Game) Scores() [2]int {
	return g.scores
}

func (g *Game) SolvedCards() []string {
	res := make([]string, len(g.cards), len(g.cards))
	for i, c := range g.solvedCards {
		res[i] = c
	}
	return res
}

func (g *Game) NumberOfCards() int {
	return len(g.cards)
}

func (g *Game) CurrentPlayer() byte {
	return g.currentPlayer
}

func (g *Game) getCard(i int) (string, error) {
	if i < 0 || i >= len(g.cards) {
		return "", fmt.Errorf("Invalid index requested")
	}
	return g.cards[i], nil
}

// TurnCard allows the given player to turn a single card if its their turn. If
// a card has been previously turned in this round and it matches, the pair is
// removed from the deck and the player's score is increased by 1.
func (g *Game) TurnCard(i int, player byte) (TurnResult, error) {
	result := TurnResult{}
	if player != g.currentPlayer {
		return result, fmt.Errorf("It's not your turn.")
	}
	result.Player = g.currentPlayer
	result.NextTurnPlayer = g.currentPlayer
	result.NewScore = g.scores[g.currentPlayer]
	c, err := g.getCard(i)
	if err != nil {
		return result, err
	}
	result.Card = c
	if c == "" {
		return result, fmt.Errorf("This card has already been removed.")
	}
	if i == g.turnedCardIndex {
		return result, fmt.Errorf("You've just turned this card.")
	}
	if g.turnedCardIndex == -1 {
		g.turnedCardIndex = i
		return result, nil
	}
	result.MatchFound = false
	result.TurnComplete = true
	if c == g.cards[g.turnedCardIndex] {
		g.solvedCards[g.turnedCardIndex] = g.cards[i]
		g.solvedCards[i] = g.cards[i]
		g.cards[i] = ""
		g.cards[g.turnedCardIndex] = ""
		result.MatchFound = true
		result.NewScore = g.increaseScore()
		g.pairsLeft--
		if g.pairsLeft == 0 {
			result.GameOver = true
			result.Winner = g.currentPlayer
		}
	} else {
		g.currentPlayer = (g.currentPlayer + 1) % 2
	}
	g.turnedCardIndex = -1
	result.NextTurnPlayer = g.currentPlayer
	return result, nil
}

func (g *Game) increaseScore() int {
	if g.currentPlayer >= 0 && g.currentPlayer <= 1 {
		g.scores[g.currentPlayer]++
	}
	return g.scores[g.currentPlayer]
}

// NewGame generates a new Game instance containing a set of cards. Due to the
// way card names are generated, there is a maximum number of different cards.
// If that is reached, an error is returned.
func NewGame(numberOfPairs int) (*Game, error) {
	if numberOfPairs > maxNumberOfPairs {
		return nil, fmt.Errorf("You are only allowed to generate a game with a maximum of %v of pairs", maxNumberOfPairs)
	}
	if numberOfPairs < 1 {
		return nil, fmt.Errorf("A game has to have at least one pair.")
	}

	cards := createSetOfCards(numberOfPairs)

	return &Game{
		cards:           cards,
		solvedCards:     make([]string, len(cards), len(cards)),
		turnedCardIndex: -1,
		currentPlayer:   0,
		pairsLeft:       len(cards) / 2,
	}, nil
}

func createSetOfCards(numberOfPairs int) []string {
	cards := make([]string, 0, 2*numberOfPairs)
	usedNames := make(map[string]struct{})

	for i := 0; i < numberOfPairs; i++ {
		for {
			name := namesgenerator.GetRandomName(0)
			if _, found := usedNames[name]; !found {
				usedNames[name] = struct{}{}
				cards = append(cards, name)
				cards = append(cards, name)
				break
			}
		}
	}

	randomizeCards(cards)
	return cards
}

func randomizeCards(cards []string) {
	perms := rand.Perm(len(cards))
	dest := make([]string, len(cards), len(cards))

	for idx, pos := range perms {
		dest[idx] = cards[pos]
	}

	for idx, val := range dest {
		cards[idx] = val
	}
}
