package matchesmodels

// SkillInfo oggetto interno ad Hero che decrive la struttura delle Skills
type SkillInfo struct {
	FileName string `json:"fileName"`
	Damage   int    `json:"damage"`
}

// Hero descrive la struttura dell'oggetto eroe
type Hero struct {
	FileName string      `json:"fileName"`
	HP       int         `json:"hp"`
	Skills   []SkillInfo `json:"skills"`
}

// AddRequest descrive il formato del JSON che un client usa per aggiungersi ad una room
type AddRequest struct {
	Room   string `json:"room"`
	Token  string `json:"token"`
	Heroes []Hero `json:"heroes"`
}

// WelcomeMessage descrive il formato del JSON che un client riceve quando entra in una room
type WelcomeMessage struct {
	Seed           int     `json:"seed"`
	IsPlayer       bool    `json:"isPlayer"`
	YourTurn       bool    `json:"yourTurn"`
	OpponentHeroes *[]Hero `json:"opponentHeroes"`
}

// ChatMessage descrive il formato del JSON che un client usa per inviare un messaggio a tutti i membri di una room
type ChatMessage struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

// ServerError descrive il formato del JSON che il server manda ad un client in seguito ad un errore
/*
0: JSON non valido per l'aggiunta in una room
1: utente già presente nella room
2: si è già loggati in una room ma il JSON non rispecchia nessuno dei formati aspettati
3: JSON valido per l'aggiunta in una room, ma token invalido
*/
type ServerError struct {
	ServerError uint `json:"serverError"`
}
