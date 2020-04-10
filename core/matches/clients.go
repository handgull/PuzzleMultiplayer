package matchesapi

import (
	matchesmodels "PuzzleMultiplayer/models/matches"
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
)

// Questo tipo rappresenta il client, interfacciandosi all'interfaccia piuttosto che all'implementazione effettiva dei canali di comunicazione
type client struct {
	*bufio.Reader
	*bufio.Writer
	wc chan string
}

// Metodo di Split usato dallo scanner (vedi *bufio.Scanner)
func scanAll(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}
	return len(data), data, nil
}

// StartClient inizializza un client e le sue logiche di funzionamento. Qui uso il channel generator pattern
func StartClient(name string, broadcastCh chan<- string, cn io.ReadWriteCloser, roomName string) (chan<- string, <-chan struct{}) {
	c := new(client)
	c.Reader = bufio.NewReader(cn)
	c.Writer = bufio.NewWriter(cn)
	c.wc = make(chan string)
	doneCh := make(chan struct{})

	// setup the reader. When the client sends a message, we will send it to the chat room
	go func() {
		scanner := bufio.NewScanner(c.Reader)
		scanner.Split(scanAll)
		for scanner.Scan() {
			msg := scanner.Text()
			myDealer(roomName, name, broadcastCh, c.wc, msg)
		}
		close(doneCh)
		cn.Close()
	}()

	c.writeMonitor()
	return c.wc, doneCh
}

func (c *client) writeMonitor() {
	go func() {
		for s := range c.wc {
			c.WriteString(s)
			c.Flush()
		}
	}()
}

// Funzione che si occupa di decodificare i messaggi ricevuti e fare diverse azioni in base al loro formato
func myDealer(roomName, name string, broadcastCh chan<- string, wc chan string, msg string) {
	chatMex := new(matchesmodels.ChatMessage)
	dec := json.NewDecoder(bytes.NewReader([]byte(msg)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&chatMex); err != nil && false {
		invalidReq := matchesmodels.ServerError{
			ServerError: 2,
		}
		bs, _ := json.Marshal(invalidReq)
		wc <- string(bs)
		log.Println("Invalid request from", name, "in", roomName, "Error:", err)
	} else {
		log.Printf("%s|%s: %s", roomName, name, msg)
		broadcastCh <- msg
	}
}
