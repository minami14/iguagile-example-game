package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/iguagile/iguagile-engine/iguagile"
)

type roomService struct {
	room       *iguagile.Room
	game       *game
	nextPlayer int
	mu         *sync.Mutex
}

func (r *roomService) gameStart() error {
	resp := responseData{
		MessageType: gameStartMessage,
		Cards:       r.game.cards,
	}

	for _, p := range r.game.players {
		resp.Player = p
		data, err := json.Marshal(&resp)
		if err != nil {
			return err
		}

		r.room.SendToClient(p.ID, 0, data)
	}

	if err := r.game.changePlayer(); err != nil {
		return err
	}

	r.game.status = inProgress
	return nil
}

func (r *roomService) Receive(senderID int, data []byte) error {
	var req requestData
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	return r.game.openCard(senderID, req.CardNumber)
}

func (r *roomService) OnRegisterClient(clientID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println(clientID, r.nextPlayer, r.game.conf.players)
	if r.nextPlayer >= r.game.conf.players {
		return fmt.Errorf("the number of participants has exceeded the maximum number")
	}

	r.game.players[r.nextPlayer] = player{ID: clientID}
	r.nextPlayer++

	if r.nextPlayer == r.game.conf.players {
		return r.gameStart()
	}

	return nil
}

func (r *roomService) OnUnregisterClient(clientID int) error {
	resp := responseData{
		MessageType: errorMessage,
		Player:      player{ID: clientID},
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	r.room.SendToAllClients(0, data)
	return nil
}

func (r *roomService) OnChangeHost(clientID int) error {
	return nil
}

func (r *roomService) Destroy() error {
	return nil
}

type roomServiceFactory struct {
	gameConfig config
}

func (r *roomServiceFactory) Create(room *iguagile.Room) (iguagile.RoomService, error) {
	return &roomService{
		room: room,
		game: newGame(room, r.gameConfig),
		mu:   new(sync.Mutex),
	}, nil
}
