package auth

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("PuzzleMultiplayer")

// Claims struct che viene encodata/decodificata
// NOTA: il tipo embedded jwt.StandardClaims fornisce i campi standard del JWT
type Claims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

// Encode restituisce il JWT corrispondente all'oggetto di tipo Claims
func (c *Claims) Encode() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(jwtKey)
}

// Decode data una stringa provo a decodificarla e a restituire un oggetto Claims
func Decode(tknStr string) (*Claims, error) {
	claims := new(Claims)
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return claims, err
	}
	if !tkn.Valid {
		return claims, fmt.Errorf("Invalid token")
	}

	return claims, nil
}
