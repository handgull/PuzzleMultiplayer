package matchesapi

import (
	matchesmodels "PuzzleMultiplayer/models/matches"
	"encoding/json"
	"io"
	"log"
	"sync"
)

// Room tipo che rappresenta una room (ovvero un match)
type Room struct {
	name      string
	Broadcast chan string
	clients   map[string]chan<- string
	Quit      chan struct{}
	*sync.RWMutex
}

// metodo che distribuisce un messaggio a tutti i client nella room
func (r *Room) broadcastMsg(msg string) {
	r.RLock()
	defer r.RUnlock()
	for _, wc := range r.clients {
		go func(wc chan<- string) {
			wc <- msg
		}(wc)
	}
}

// Funzione usata per non ripetere codice tra RemoveClient e CloseRoom
func removeClientNoMutex(r *Room, name string) {
	log.Println("Removing client", name)
	delete(r.clients, name)
}

// RemoveClientSync rimuove un client dalla room. (blocking call)
func (r *Room) RemoveClientSync(name string) {
	r.Lock()
	defer r.Unlock()
	removeClientNoMutex(r, name)
}

// CloseChatRoomSync chiude una room. (blocking call)
func (r *Room) CloseChatRoomSync() {
	r.Lock()
	defer r.Unlock()
	log.Println("Closing room", r.name)
	for name := range r.clients {
		removeClientNoMutex(r, name)
	}
	close(r.Broadcast)
	r.Quit <- struct{}{}
	close(r.Quit)
}

// Run fa partire in diverse goroutines le logiche di una room
func (r *Room) Run() {
	log.Println("Starting chat room", r.name)
	// Ogni messaggio ricevuto sul channel di broadcast va distribuito a tutti i client
	go func() {
		for msg := range r.Broadcast {
			r.broadcastMsg(msg)
		}
	}()

	// Se ricevo un segnale sul channel Quit allora la room va terminata
	go func() {
		<-r.Quit
		r.CloseChatRoomSync()
	}()
}

// AddClient aggiunge un client alla room
func (r *Room) AddClient(c io.ReadWriteCloser, clientname string, seed int, isPlayer, yourTurn bool, opponentHeroes *[]matchesmodels.Hero) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.clients[clientname]; ok {
		log.Printf("Client %s already exist in chat room %s! Closing the connection", clientname, r.name)
		invalidReq := matchesmodels.ServerError{
			ServerError: 1,
		}
		bs, _ := json.Marshal(invalidReq)
		c.Write(bs)
		c.Close()
		return
	}
	log.Printf("Adding client %s in %s, clients number BEFORE the insert: %v", clientname, r.name, len(r.clients))
	wc, done := StartClient(clientname, r.Broadcast, c, r.name)
	r.clients[clientname] = wc
	welcomeMsg := matchesmodels.WelcomeMessage{
		Seed:           seed,
		IsPlayer:       isPlayer,
		YourTurn:       yourTurn,
		OpponentHeroes: opponentHeroes,
	}
	bs, _ := json.Marshal(welcomeMsg)
	c.Write(bs)

	// Se ricevo un segnale sul channel done allora devo rimuovere il client
	go func() {
		<-done
		r.RemoveClientSync(clientname)
	}()
}

// CreateRoom inizializza ed avvia le logiche di una room
func CreateRoom(rname string) *Room {
	r := &Room{
		name:      rname,
		Broadcast: make(chan string),
		RWMutex:   new(sync.RWMutex),
		clients:   make(map[string]chan<- string),
		Quit:      make(chan struct{}),
	}
	r.Run()
	return r
}
