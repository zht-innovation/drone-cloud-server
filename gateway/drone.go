package gateway

import (
	"log"
	"net/http"
	"encoding/json"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}

func droneDataHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket conn build error: %v", err)
	}

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Websocket read data error: %v", err)
		}
		
		if messageType == 1 { // Byte stream
			var droneData DroneData
			if err := json.Unmarshal(p, &droneData); err != nil {
				log.Fatalf("Json unmarshal data error: %v", err)
			}
		}
	}
}