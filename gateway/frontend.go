// v0: Receive messages from Redis and send them to the frontend via WebSocket
package gateway

import (
	"context"
	"net/http"
	"sync"
	S "zhtcloud/gateway/shared"
	"zhtcloud/utils"
	"zhtcloud/utils/logger"

	ws "github.com/gorilla/websocket"
)

var mu sync.Mutex

func runSendRunningStatus(ctx context.Context, c chan string, conn *ws.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			status := <-c
			mu.Lock()
			err := conn.WriteMessage(ws.TextMessage, []byte(status))
			mu.Unlock()
			if err != nil {
				if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseNormalClosure) {
					logger.Error("Write msg to web: %v", err)
				} else {
					return
				}
			}
		}
	}
}

func frontendTransfer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(S.Ctx)
	conn, err := S.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Websocket conn build error: %v", err)
	}

	defer conn.Close()
	defer cancel()

	redisClient := utils.GetRedisConn(ctx)

	droneInfoChan := make(chan string)
	runningStatusChan := make(chan string)

	go func() {
		if err := redisClient.SubChannel(ctx, "drone_info", droneInfoChan); err != nil {
			logger.Error("Error subscribing to channel: %v", err)
		}
	}()

	go func() {
		if err := redisClient.SubChannel(ctx, "running_status", runningStatusChan); err != nil {
			logger.Error("Error subscribing to channel: %v", err)
		}
	}()

	go runSendRunningStatus(ctx, runningStatusChan, conn)

	for {
		data := <-droneInfoChan
		mu.Lock()
		err := conn.WriteMessage(ws.TextMessage, []byte(data))
		mu.Unlock()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseNormalClosure) {
				logger.Error("Write msg to web: %v", err)
			} else {
				return
			}
		}
	}
}
