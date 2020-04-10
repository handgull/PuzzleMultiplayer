package matchesapi

import (
	matchesmodels "PuzzleMultiplayer/models/matches"
	"math/rand"
	"sync"
	"time"
)

type matchMakerData struct {
	seed           int
	room           string
	yourTurn       bool
	OpponentHeroes *[]matchesmodels.Hero
}

type opponent struct {
	name   string
	heroes *[]matchesmodels.Hero
	dataCh chan matchMakerData
}

// Struttura che contiene i giocatori in attesa che siano smistati in delle room
type waitingRoom struct {
	opponents []*opponent
	Quit      chan struct{}
	*sync.RWMutex
}

func (w *waitingRoom) pair(name string, heroes *[]matchesmodels.Hero, dc chan matchMakerData) {
	w.Lock()
	defer w.Unlock()
	if len(w.opponents) > 0 {
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		data := matchMakerData{
			seed:           r.Intn(100000),
			room:           w.opponents[0].name + "VS" + name,
			yourTurn:       false,
			OpponentHeroes: w.opponents[0].heroes,
		}
		go func(d matchMakerData) {
			dc <- d
		}(data)
		go func(d matchMakerData) {
			d.yourTurn, d.OpponentHeroes = true, heroes
			w.opponents[0].dataCh <- d
			w.opponents = w.opponents[1:]
		}(data)
	} else {
		w.opponents = append(w.opponents, &opponent{
			name:   name,
			dataCh: dc,
			heroes: heroes,
		})
	}

}

// inizializza ed avvia le logiche di una waitingRoom
func createWaitingRoom() *waitingRoom {
	wr := &waitingRoom{
		RWMutex: new(sync.RWMutex),
		Quit:    make(chan struct{}),
	}
	// wr.Run()
	return wr
}
