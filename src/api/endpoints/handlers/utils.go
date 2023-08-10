package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"strconv"
)

func convertToInt(value string, paramName string) (int, error) {
	if value == "" {
		err := errors.New(fmt.Sprintf("Parâmetro não informado: %s", paramName))
		log.Error("Requisição mal formulada: ", err.Error())
		return -1, err
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		errorReturned := errors.New(fmt.Sprintf("Parâmetro inválido: %s", paramName))
		log.Errorf("Requisição mal formulada: %s - %s", errorReturned, err.Error())
		return -1, errorReturned
	}

	return intValue, nil
}
