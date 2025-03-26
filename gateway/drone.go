package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"zhtcloud/utils"

	ws "github.com/gorilla/websocket"
)

func droneDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket conn build error: %v", err)
	}

	defer conn.Close()
	defer cancel()

	redisClient := utils.GetRedisConn(ctx)

	msgChan := make(chan string)

	// listen coordinates transfer
	go func() {
		if err := redisClient.SubChannel(ctx, COORS, msgChan); err != nil {
			log.Printf("Error subscribing to 'coords' channel: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				coors := <-msgChan
				err := conn.WriteMessage(ws.TextMessage, []byte(coors))
				if err != nil {
					log.Printf("Send coordinates to drones: %v", err)
				}
			}
		}
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Websocket read data error: %v", err)
		}

		if messageType == 1 { // Byte stream
			var droneData DroneData
			if err := json.Unmarshal(p, &droneData); err != nil {
				log.Printf("Json unmarshal data error: %v", err)
			}
			redisClient.PubChannel(ctx, "tmp", string(p))
		}
	}
}
