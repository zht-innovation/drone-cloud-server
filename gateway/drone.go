// Read message from drones and send them to Redis
package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"zhtcloud/utils"
	"zhtcloud/utils/logger"

	ws "github.com/gorilla/websocket"
)

func droneDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Websocket conn build error: %v", err)
	}

	defer conn.Close()
	defer cancel()

	redisClient := utils.GetRedisConn(ctx)

	msgChan := make(chan string)

	// listen coordinates transfer
	go func() {
		if err := redisClient.SubChannel(ctx, COORS, msgChan); err != nil {
			logger.Error("Error subscribing to 'coords' channel: %v", err)
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
					if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseNormalClosure) {
						logger.Error("Send coordinates to drones: %v", err)
					} else {
						return
					}
				}
			}
		}
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure,
				ws.CloseNormalClosure, ws.CloseNoStatusReceived) { // NoStatusReceived: websocket 1005
				logger.Error("Websocket read data: %v", err)
			} else {
				return
			}
		}

		if messageType == 1 { // Byte stream
			decoder := json.NewDecoder(bytes.NewReader(p))

			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				logger.Fatal("JSON decode error: %v", err)
			}

			typ := uint8(data["TYPE"].(float64)) // from interface{} default float64
			var chanName string

			switch typ {
			case DRONE_INFO:
				chanName = "drone_info"
			case RUNNING_STATUS:
				chanName = "running_status"
			}

			redisClient.PubChannel(ctx, chanName, string(p))
		}
	}
}
