package httpcore

import (
	"MinerGame/repositor"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type BuyMinerRequest struct {
	Type string `json:"type"`
}

type BuyMinerResponse struct {
	Miner repositor.MinerInfo `json:"miner"`
}

type BuyAnyDTO struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string    `json:"error"`
	Time  time.Time `json:"time"`
}

func NewErrorDTO(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
		Time:  time.Now(),
	}
}

func WriteError(w http.ResponseWriter, err error, status int) {
	errDTO := NewErrorDTO(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errDTO)
}
func WriteResponse(w http.ResponseWriter, body any, status int) {
	response, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error writing response:", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(response); err != nil {
		fmt.Println("Error to Response" + err.Error())
	}
}

func SelectError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repositor.ErrorNotEnoughCoal):
		WriteError(w, err, http.StatusConflict)

	case errors.Is(err, repositor.ErrorTypeNotFound):
		WriteError(w, err, http.StatusNotFound)

	case errors.Is(err, repositor.ErrorEquipmentAlreadyBuy):
		WriteError(w, err, http.StatusConflict)

	case errors.Is(err, repositor.ErrorEquipmentNotFound):
		WriteError(w, err, http.StatusNotFound)

	default:
		WriteError(w, err, http.StatusInternalServerError)
	}
}
