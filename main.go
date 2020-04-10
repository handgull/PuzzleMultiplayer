package main

import (
	matchescontroller "PuzzleMultiplayer/delivery/websockets"
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	wsAddr = flag.String("ws", "localhost:9000", "Address for the websocket matches server to listen on")
)

func main() {
	flag.Parse()
	// Faccio partire l'api che gestisce i match su websocket
	go func() {
		if err := matchescontroller.RunWSMatches("/matches", *wsAddr, 8192, 8192, 4096); err != nil {
			log.Panic(err)
		}
	}()

	// Blocco l'esecuzione in attesa di un input da tastiera (qui in futuro ci sarà la serve dell'api rest)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
