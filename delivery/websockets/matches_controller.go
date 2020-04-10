package matchescontroller

import (
	matchesapi "PuzzleMultiplayer/core/matches"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Tipo che uso per wrappare la connessione e far rispettare l'interfaccia ReadWriteCloser, per avere le logiche dell'api slegate dal websocket
type websocketConnWrapper struct {
	*websocket.Conn
	msgType int
}

// Tipo che uso nell'handler per esporre effettivamente l'api sul socket
type wsMatchesHandler struct {
	*websocket.Upgrader
	maxMessageSize int64
	*matchesapi.MatchesAPI
}

/***
Implementazione dei metodi Read Write Close per rispettare l'interfaccia ReadWriteCloser descritta nelle doc degli standard packages
***/
func (wsWrapper *websocketConnWrapper) Read(bs []byte) (int, error) {
	t, r, err := wsWrapper.NextReader()
	wsWrapper.msgType = t
	if err != nil {
		log.Println("Websocket err on Read method:", err)
		return 0, err
	}
	return r.Read(bs)
}

func (wsWrapper *websocketConnWrapper) Write(bs []byte) (int, error) {
	err := wsWrapper.WriteMessage(wsWrapper.msgType, bs)
	if err != nil {
		log.Println("Websocket err on Write method:", err)
		return 0, err
	}
	return len(bs), nil
}

// NOTA: Il metodo Close è già implementato da *websocket.Conn

// RunWSMatchesWithExistingAPI fa partire su un websocket un match server usando un API esistente
func RunWSMatchesWithExistingAPI(url, address string, Rbsize, Wbsize int, maxSize int64, matches *matchesapi.MatchesAPI) error {
	handler := &wsMatchesHandler{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  Rbsize,
			WriteBufferSize: Wbsize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		maxMessageSize: maxSize,
		MatchesAPI:     matches,
	}
	http.HandleFunc(url, handler.wshandler)
	return http.ListenAndServe(address, nil)
}

// RunWSMatches fa partire su un websocket un match server
func RunWSMatches(url, address string, Rbsize, Wbsize int, maxSize int64) error {
	return RunWSMatchesWithExistingAPI(url, address, Rbsize, Wbsize, maxSize, matchesapi.New())
}

// Funzione esposta nell'endpoint
func (wh *wsMatchesHandler) wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wh.Upgrade(w, r, nil)
	if err != nil {
		log.Panic(w, "Error occured while trying to upgrade websocket", err)
	}

	conn.SetReadLimit(wh.maxMessageSize)

	wsWrapper := &websocketConnWrapper{Conn: conn}
	wh.AddClient(wsWrapper)
}
