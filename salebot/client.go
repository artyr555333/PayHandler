package salebot

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

type Message struct { //salebot
	MessageID string `json:"message_id"`
	ClientID  string `json:"client_id"`
}

func Sale(clientID string) error {

	message := Message{MessageID: "34441949", ClientID: "550512674"}
	bytes, _ := json.Marshal(message)
	response, err := http.Post(os.Getenv("API_SALEBOT_URL")+"message", "application/json", strings.NewReader(string(bytes)))

	if err != nil {
		return err
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {

		return errors.New(string(bodyBytes))
	}
	return nil
}

func SaleAsync(resultChan chan<- error, clientID string) {
	go func() {
		err := Sale(clientID)
		resultChan <- err
	}()
}
