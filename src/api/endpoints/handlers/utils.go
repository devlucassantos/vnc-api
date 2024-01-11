package handlers

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"strconv"
	"time"
)

func convertToInt(value string, paramName string) (int, error) {
	if value == "" {
		err := errors.New(fmt.Sprintf("Parâmetro não informado: %s", paramName))
		log.Error("Requisição mal formulada: ", err.Error())
		return -1, err
	}
	intValue, err := strconv.Atoi(value)
	if err != nil || intValue <= 0 {
		errorReturned := errors.New(fmt.Sprintf("Parâmetro inválido: %s", paramName))
		log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorReturned, value)
		return -1, errorReturned
	}

	return intValue, nil
}

func convertToUuid(value string, paramName string) (uuid.UUID, error) {
	if value == "" {
		err := errors.New(fmt.Sprintf("Parâmetro não informado: %s", paramName))
		log.Error("Requisição mal formulada: ", err.Error())
		return uuid.UUID{}, err
	}

	uuidValue, err := uuid.Parse(value)
	if err != nil || uuidValue.ID() == 0 {
		errorReturned := errors.New(fmt.Sprintf("Parâmetro inválido: %s", paramName))
		log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorReturned, value)
		return uuid.UUID{}, errorReturned
	}

	return uuidValue, nil
}

func convertToTime(value string, paramName string) (time.Time, error) {
	if value == "" {
		err := errors.New(fmt.Sprintf("Parâmetro não informado: %s", paramName))
		log.Error("Requisição mal formulada: ", err.Error())
		return time.Time{}, err
	}

	timeValue, err := time.Parse("2006-01-02", value)
	if err != nil {
		errorReturned := errors.New(fmt.Sprintf("Parâmetro inválido: %s", paramName))
		log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorReturned, value)
		return time.Time{}, errorReturned
	}

	return timeValue, nil
}
