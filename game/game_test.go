package game

import "testing"

func TestNewGame(t *testing.T) {
	var g PairsGame
	var err error
	_, err = NewGame(0)
	if err == nil {
		t.Error("Creating a game with 0 pairs should be impossible")
	}

	_, err = NewGame(maxNumberOfPairs + 1)
	if err == nil {
		t.Error("Creating a game with more than max number of pairs should be impossible")
	}

	g, err = NewGame(10)
	if g == nil {
		t.Error("Creating a game with 2 pairs should be possible")
	}

	if err != nil {
		t.Errorf("No error should be generated for a valid pair number. Error received: %s", err.Error())
	}
}

func TestGetCard(t *testing.T) {
	g, _ := NewGame(10)

	_, err := g.getCard(-1)
	if err == nil {
		t.Errorf("Requesting a card with an index < 0 should have returned an error")
	}

	_, err = g.getCard(20)
	if err == nil {
		t.Errorf("Requesting a card with an index > maximumNumberOfCards should have returned an error")
	}

	c, err := g.getCard(5)
	if c == "" {
		t.Errorf("A valid index should have returned a non-empty string")
	}
	if err != nil {
		t.Errorf("A valid index should not have returned an error. Error received: %s", err.Error())
	}
}

func TestTurnCard(t *testing.T) {
	g, _ := NewGame(2)
	g.cards = []string{"a", "b", "a", "b"}

	res, err := g.TurnCard(0, 0)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if res.NewScore != 0 {
		t.Errorf("Score should still have been 0")
	}
	if res.Card == "" {
		t.Errorf("There should have been a card")
	}
	if res.MatchFound {
		t.Errorf("Finding a match on the first turn is ... hard")
	}
	if res.TurnComplete {
		t.Errorf("The round isn't complete after the first card-turn")
	}
	if res.Player != 0 {
		t.Errorf("The player should be the current one")
	}
	if res.NextTurnPlayer != 0 {
		t.Errorf("The next turn should still be done by the current player")
	}

	// Let's give the other player a go to test some additional flags
	res, err = g.TurnCard(1, 0)
	if res.NextTurnPlayer != 1 {
		t.Errorf("It should be the other player's turn now ...")
	}

	_, err = g.TurnCard(0, 1)
	if err != nil {
		t.Errorf("Unexpected error on first turn of player 2: %s", err.Error())
	}
	res, err = g.TurnCard(2, 1)
	if !res.TurnComplete {
		t.Errorf("The turn should be over now")
	}
	if !res.MatchFound {
		t.Errorf("This should have been a match")
	}
	if res.NewScore != 1 {
		t.Errorf("This player should have scored")
	}
	if res.NextTurnPlayer != 1 {
		t.Errorf("Someone who finds a match should have another round")
	}
	g.TurnCard(1, 1)
	res, _ = g.TurnCard(3, 1)
	if !res.GameOver {
		t.Errorf("The game should be over now")
	}
	if res.Winner != 1 {
		t.Errorf("Player 2 should be the winner")
	}
}
