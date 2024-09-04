package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
	"vnc-api/api/endpoints/dto/response"
)

func ConvertFromStringToInt(value string, paramName string) (int, *response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parâmetro não informado: %s", paramName)
		log.Error("Requisição mal formulada: ", errorMessage)
		return -1, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil || intValue <= 0 {
		errorMessage := fmt.Sprintf("Parâmetro inválido: %s", paramName)
		log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, value)
		return -1, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	return intValue, nil
}

func ConvertFromStringToUuid(value string, paramName string) (uuid.UUID, *response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parâmetro não informado: %s", paramName)
		log.Error("Requisição mal formulada: ", errorMessage)
		return uuid.Nil, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	uuidValue, err := uuid.Parse(value)
	if err != nil || uuidValue.ID() == 0 {
		errorMessage := fmt.Sprintf("Parâmetro inválido: %s", paramName)
		log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, value)
		return uuid.Nil, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	return uuidValue, nil
}

func ConvertFromStringToTime(value string, paramName string) (time.Time, *response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parâmetro não informado: %s", paramName)
		log.Error("Requisição mal formulada: ", errorMessage)
		return time.Time{}, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	timeValue, err := time.Parse("2006-01-02", value)
	if err != nil {
		errorMessage := fmt.Sprintf("Parâmetro inválido: %s", paramName)
		log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, value)
		return time.Time{}, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	return timeValue, nil
}
