package matchesapi

import (
	"PuzzleMultiplayer/auth"
	matchesmodels "PuzzleMultiplayer/models/matches"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// innerAddClient descrive il formato interno per aggiungere un client ad una room
type innerAddClient struct {
	room     string
	username string
}

// MatchesAPI tipo che gestisce le varie room
type MatchesAPI struct {
	rooms map[string]*Room
	*waitingRoom
	*sync.RWMutex
}

// New inizializza e fa partire le logiche di funzionamento dell' api dei matches
func New() *MatchesAPI {
	api := &MatchesAPI{
		rooms:       make(map[string]*Room),
		waitingRoom: createWaitingRoom(),
		RWMutex:     new(sync.RWMutex),
	}

	go func() {
		// Se l'applicazione riceve dal kernel SIGINT o SIGTERM invio il segnale sul channel ch (di conseguenza chiudo tutte le room)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Println("Closing connection")
		api.RLock()
		defer api.RUnlock()
		// Questo ciclo manda il segnale vuoto sul canale Quit di ogni room e poi aspetta una risposta sul medesimo canale per capire che è stata terminata con successo
		for _, r := range api.rooms {
			r.Quit <- struct{}{}
			<-r.Quit
		}
		os.Exit(0)
	}()

	return api
}

// AddClient adds a new client to the matches api. Expects a JSON file
func (mAPI *MatchesAPI) AddClient(c io.ReadWriteCloser) {
	req := new(matchesmodels.AddRequest)
	dec := json.NewDecoder(c)
	dec.DisallowUnknownFields()
	if err := dec.Decode(req); err != nil {
		invalidReq := matchesmodels.ServerError{
			ServerError: 0,
		}
		bs, _ := json.Marshal(invalidReq)
		c.Write(bs)
		log.Println("Could not create chat room:", err)
		c.Close()
	} else {
		token, err := auth.Decode(req.Token)
		if err != nil {
			invalidReq := matchesmodels.ServerError{
				ServerError: 3,
			}
			bs, _ := json.Marshal(invalidReq)
			c.Write(bs)
			log.Println("Could not create chat room, invalid JWT:", err)
			c.Close()
			return
		}
		uname := strconv.FormatInt(int64(token.ID), 10)
		isPlayer := req.Room == ""
		yourTurn := false
		seed := -1
		opponentHeroes := new([]matchesmodels.Hero)
		if isPlayer {
			mmakerDataCh := make(chan matchMakerData)
			mAPI.waitingRoom.pair(uname, &req.Heroes, mmakerDataCh)
			matchmakingRes := <-mmakerDataCh
			seed, req.Room, yourTurn, opponentHeroes = matchmakingRes.seed, matchmakingRes.room, matchmakingRes.yourTurn, matchmakingRes.OpponentHeroes
		}
		// TODO: soluzione temporanea, qui metto brutalmente l'id dal token
		internalReq := &innerAddClient{
			room:     req.Room,
			username: uname,
		}

		mAPI.handleClient(internalReq, c, seed, isPlayer, yourTurn, opponentHeroes)
	}
}

// Metodo che gestisce le richieste di aggiungersi ad una room che arrivano da un client
func (mAPI *MatchesAPI) handleClient(req *innerAddClient, c io.ReadWriteCloser, seed int, isPLayer, yourTurn bool, opponentHeroes *[]matchesmodels.Hero) {
	mAPI.Lock()
	defer mAPI.Unlock()
	r, ok := mAPI.rooms[req.room]
	if !ok {
		r = CreateRoom(req.room)
	}
	r.AddClient(c, req.username, seed, isPLayer, yourTurn, opponentHeroes)
	mAPI.rooms[req.room] = r
}
