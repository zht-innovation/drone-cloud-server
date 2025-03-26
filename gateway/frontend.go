package gateway

import (
	"log"
	"net/http"
	"zhtcloud/utils"

	ws "github.com/gorilla/websocket"
)

func frontendTransfer(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket conn build error: %v", err)
	}

	defer conn.Close()

	redisClient := utils.GetRedisConn(ctx)

	msgChan := make(chan string)

	go func() {
		if err := redisClient.SubChannel(ctx, "tmp", msgChan); err != nil {
			log.Printf("Error subscribing to channel: %v", err)
		}
	}()

	for {
		data := <-msgChan
		err := conn.WriteMessage(ws.TextMessage, []byte(data))
		if err != nil {
			log.Printf("Write msg to web: %v", err)
		}
	}
}
