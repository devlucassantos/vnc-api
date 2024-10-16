package utils

import (
	"crypto/rand"
	"github.com/labstack/gommon/log"
	"math/big"
)

func GenerateUserActivationCode() (string, error) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	activationCode := make([]byte, length)
	for index := range activationCode {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Error("Error generating user account activation code: ", err.Error())
			return "", err
		}

		activationCode[index] = charset[randomInt.Int64()]
	}

	return string(activationCode), nil
}
