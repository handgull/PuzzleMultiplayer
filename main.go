package main

import (
	matchescontroller "PuzzleMultiplayer/delivery/websockets"
	"flag"
	"log"
)

var (
	wsAddr = flag.String("ws", "localhost:9000", "Address for the websocket matches server to listen on")
)

func main() {
	flag.Parse()
	// Faccio partire l'api che gestisce i match su websocket
	if err := matchescontroller.RunWSMatches("/matches", *wsAddr, 8192, 8192, 4096); err != nil {
		log.Panic(err)
	}
}
