package waygate

import (
	"crypto/rand"
	"math/big"
)

type Waygate struct {
	Domains     []string `json:"domains"`
	Description string   `json:"description"`
}

type TokenData struct {
	WaygateId string `json:"waygate_id"`
}

func GenRandomCode() (string, error) {

	const chars string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	id := ""
	for i := 0; i < 32; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		id += string(chars[randIndex.Int64()])
	}
	return id, nil
}
