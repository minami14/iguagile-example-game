package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/iguagile/iguagile-engine/iguagile"
)

type player struct {
	ID    int `json:"id"`
	Score int `json:"score"`
}

type card struct {
	ID       int `json:"id"`
	Number   int `json:"number"`
	acquired bool
}

func generateCards(pairs int) []*card {
	cards := make([]*card, pairs*2)
	for i := range cards {
		cards[i] = &card{ID: i / 2}
	}

	for i := len(cards) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}

	for i := 0; i < len(cards); i++ {
		cards[i].Number = i
	}

	return cards
}

type config struct {
	cardPairs int
	players   int
}

// Game status
const (
	inPreparation = iota
	inProgress
	finish
)

type game struct {
	room          *iguagile.Room
	cards         []*card
	openedCard    *card
	players       []player
	conf          config
	status        int
	currentPlayer int
	acquired      int
	mu            *sync.Mutex
}

func newGame(room *iguagile.Room, conf config) *game {
	return &game{
		room:          room,
		cards:         generateCards(conf.cardPairs),
		players:       make([]player, conf.players),
		conf:          conf,
		status:        inPreparation,
		currentPlayer: -1,
		mu:            new(sync.Mutex),
	}
}

func (g *game) restart() {
	g.cards = generateCards(g.conf.cardPairs)
	for _, p := range g.players {
		p.Score = 0
	}
}

func (g *game) openCard(senderID, cardNumber int) error {
	if g.status != inProgress {
		return fmt.Errorf("the game is not in progress")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	currentPlayerID := g.players[g.currentPlayer].ID
	if senderID != currentPlayerID {
		return fmt.Errorf("invalid player. current: %v, sender: %v", currentPlayerID, senderID)
	}

	if g.cards[cardNumber].acquired {
		return fmt.Errorf("the card has already been acquired")
	}

	if g.openedCard == nil {
		g.openedCard = g.cards[cardNumber]
		resp := responseData{
			MessageType: openCardMessage,
			Card:        g.openedCard,
			Player:      g.players[g.currentPlayer],
		}

		data, err := json.Marshal(&resp)
		if err != nil {
			return err
		}

		g.room.SendToAllClients(0, data)
		return nil
	}

	if g.cards[cardNumber].Number == g.openedCard.Number {
		return nil
	}

	resp := responseData{
		MessageType: openCardMessage,
		Card:        g.cards[cardNumber],
		Player:      g.players[g.currentPlayer],
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	g.room.SendToAllClients(0, data)

	time.Sleep(time.Second * 1)

	if g.cards[cardNumber].ID == g.openedCard.ID {
		return g.getCard(currentPlayerID, cardNumber)
	}

	if err := g.closeCard(); err != nil {
		return err
	}

	if err := g.changePlayer(); err != nil {
		return err
	}

	return nil
}

func (g *game) getCard(playerID, cardNumber int) error {
	g.openedCard.acquired = true
	g.cards[cardNumber].acquired = true
	g.players[playerID].Score++
	resp := responseData{
		MessageType: getCardMessage,
		Cards:       []*card{g.openedCard, g.cards[cardNumber]},
		Player:      g.players[playerID],
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	g.room.SendToAllClients(0, data)
	g.openedCard = nil
	g.acquired++
	if g.acquired == g.conf.cardPairs {
		return g.finish()
	}

	return nil
}

func (g *game) closeCard() error {
	resp := responseData{
		MessageType: closeCardMessage,
		Card:        g.openedCard,
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	g.openedCard = nil
	g.room.SendToAllClients(0, data)
	return nil
}

func (g *game) changePlayer() error {
	g.currentPlayer++
	if g.currentPlayer == g.conf.players {
		g.currentPlayer = 0
	}

	resp := responseData{
		MessageType: changePlayerMessage,
		Player:      g.players[g.currentPlayer],
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	g.room.SendToAllClients(0, data)

	return nil
}

func (g *game) finish() error {
	g.status = finish
	winner := g.players[0]
	for _, p := range g.players {
		if winner.Score < p.Score {
			winner = p
		}
	}

	resp := responseData{
		MessageType: finishGameMessage,
		Player:      winner,
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	g.room.SendToAllClients(0, data)

	return nil
}
