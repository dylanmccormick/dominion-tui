package main

import (
	"fmt"

	"github.com/dylanmccormick/dominion-tui/internal/cards"
	"github.com/dylanmccormick/dominion-tui/server"
)

type Action string
type Card struct {
	Name string
	Image []byte //TODO
	Action Action
	Cost int
}

type OrderedDeck struct {
	Cards []Card
}

type UnorderedDeck struct {
	CardMap map[string]int // Map [cardName]amount
}


func main() {

	commonDeck := buildStartingDeck()
	for k, v := range commonDeck.CardMap {
		fmt.Printf("Card Info: %s\nCount: %d\n", cards.CardDict[k].String(), v)
	}

	s := server.Init("42069")
	s.Serve()

	

}


func (od *OrderedDeck) peek(n int) {
	// check the top n cards of the deck and return them 
}

func buildStartingDeck() *UnorderedDeck {
	cm := map[string]int{
		"cellar": 10,
		"market": 10,
		"merchant": 10,
		// "militia": 10,
		// "mine": 10,
		// "moat": 10, 
		// "remodel": 10,
		// "smithy": 10,
		// "village": 10,
		// "workshop": 10,
	}

	return &UnorderedDeck {
		CardMap: cm,
	}

}

func GetCardFromString(s string) {

}

