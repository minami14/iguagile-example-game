package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/iguagile/iguagile-engine/iguagile"
)

func main() {
	factory := &roomServiceFactory{
		gameConfig: config{
			cardPairs: 10,
			players:   2,
		},
	}

	store, err := iguagile.NewRedis(os.Getenv("REDIS_HOST"))
	if err != nil {
		log.Fatal(err)
	}

	address := os.Getenv("ROOM_HOST")
	server, err := iguagile.NewRoomServer(factory, store, address)
	if err != nil {
		log.Fatal(err)
	}
	server.ServerUpdateDuration = time.Second * 10

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Run(listener, port))
}
