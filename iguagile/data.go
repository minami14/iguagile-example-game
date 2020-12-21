package main

// Message types
const (
	gameStartMessage    = "start"
	openCardMessage     = "open"
	closeCardMessage    = "close"
	getCardMessage      = "get"
	changePlayerMessage = "change"
	finishGameMessage   = "finish"
	errorMessage        = "error"
)

type requestData struct {
	CardNumber int `json:"card_number"`
}

type responseData struct {
	MessageType string  `json:"message_type"`
	Cards       []*card `json:"cards,omitempty"`
	Card        *card   `json:"card,omitempty"`
	Player      player  `json:"player,omitempty"`
}
