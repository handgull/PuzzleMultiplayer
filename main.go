package main

import (
	matchescontroller "PuzzleMultiplayer/delivery/websockets"
	"flag"
	"log"
)

var (
	wsAddr = flag.String("ws", ":9000", "Address for the websocket matches server to listen on")
)

func main() {
	flag.Parse()
	// Faccio partire l'api che gestisce i match su websocket
	if err := matchescontroller.RunWSMatches("/ws", *wsAddr, 8192, 8192, 4096); err != nil {
		log.Panic(err)
	}
}
