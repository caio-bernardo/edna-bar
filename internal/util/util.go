package util

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	RequestTimeout = 2 * time.Second
	ErrInvalidID  = errors.New("invalid id parameter")
)
/// Escreve uma reposta com o corpo em JSON com o status passado
func WriteJSON(w http.ResponseWriter, status int, v any ) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err = w.Write(res); err != nil {
		return err
	}
	return nil
}

/// Lê o corpo (em json) da requisição, decodifica e armazena no destino
func ReadJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func GetIDParam(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidID
	}
	return id, nil
}

/// Escreve uma mensagem de error com o status passado, o corpo da mensagem será em JSON
func ErrorJSON(w http.ResponseWriter, msg string, status int) {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res, err := json.Marshal(map[string]string{"message": msg})
	// Impossivel
	if err != nil {
		log.Printf("Error ao criar mensagem em json: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
