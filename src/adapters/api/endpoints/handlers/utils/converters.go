package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
	"vnc-api/adapters/api/endpoints/dto/response"
)

func ConvertFromStringToBool(value string, parameter string, parameterDescription string) (bool, *response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parameter not provided: %s (%s)", parameterDescription, parameter)
		log.Warn("Badly formatted request: ", errorMessage)
		return false, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid parameter: %s (%s)", parameterDescription, parameter)
		log.Warnf("Badly formatted request: %s (Value: %s)", errorMessage, value)
		return false, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	return boolValue, nil
}

func ConvertFromStringToInt(value string, parameter string, parameterDescription string) (int, *response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parameter not provided: %s (%s)", parameterDescription, parameter)
		log.Warn("Badly formatted request: ", errorMessage)
		return -1, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid parameter: %s (%s)", parameterDescription, parameter)
		log.Warnf("Badly formatted request: %s (Value: %s)", errorMessage, value)
		return -1, response.NewHttpError(http.StatusBadRequest, errorMessage)
	} else if intValue <= 0 {
		errorMessage := fmt.Sprintf("Invalid parameter: %s (%s) must be greater than 0", parameterDescription,
			parameter)
		log.Warnf("Parameter out of allowed range: %s (Value: %s)", errorMessage, value)
		return -1, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage)
	}

	return intValue, nil
}

func ConvertFromStringToTime(value string, parameter string, parameterDescription string) (time.Time,
	*response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parameter not provided: %s (%s)", parameterDescription, parameter)
		log.Warn("Badly formatted request: ", errorMessage)
		return time.Time{}, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	timeValue, err := time.Parse("2006-01-02", value)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid parameter: %s (%s)", parameterDescription, parameter)
		log.Warnf("Badly formatted request: %s (Value: %s)", errorMessage, value)
		return time.Time{}, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	return timeValue, nil
}

func ConvertFromStringToUuid(value string, parameter string, parameterDescription string) (uuid.UUID,
	*response.HttpError) {
	if value == "" {
		errorMessage := fmt.Sprintf("Parameter not provided: %s (%s)", parameterDescription, parameter)
		log.Warn("Badly formatted request: ", errorMessage)
		return uuid.Nil, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	uuidValue, err := uuid.Parse(value)
	if err != nil || uuidValue.ID() == 0 {
		errorMessage := fmt.Sprintf("Invalid parameter: %s (%s)", parameterDescription, parameter)
		log.Warnf("Badly formatted request: %s (Value: %s)", errorMessage, value)
		return uuid.Nil, response.NewHttpError(http.StatusBadRequest, errorMessage)
	}

	return uuidValue, nil
}
