package salebot

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type Message struct { //salebot
	MessageID string `json:"message_id"`
	ClientID  string `json:"client_id"`
}

func Sale(clientID string) error {

	message := Message{MessageID: "34441949", ClientID: "550512674"}
	bytes, _ := json.Marshal(message)
	response, err := http.Post("https://chatter.salebot.pro/api/a3eb382ea39c4b25bd336eec08aca028/message", "application/json", strings.NewReader(string(bytes)))
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {

		return errors.New(string(bodyBytes))
	}
	return err
}

func SaleAsync(resultChan chan<- error, clientID string) {
	go func() {
		log.Println("функция SaleAsync")
		err := Sale(clientID)
		if err != nil {
			log.Println("Ошибка записана в Канал")
			resultChan <- err
		}
	}()
}
