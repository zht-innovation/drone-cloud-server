package main

import (
	"encoding/json"
	"io"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func _setupClient(pwd string) bool {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ws://110.42.101.86:8083/mqtt")
	opts.SetClientID("websocket")
	opts.SetUsername(pwd)
	opts.SetPassword("")
	cli := mqtt.NewClient(opts)

	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return true
}

func main() {
	rsp, _ := http.Get("http://localhost:32223/auth?mac=00-11-22-33-44-55")
	if rsp.StatusCode == 200 {
		buf, _ := io.ReadAll(rsp.Body)
		res := response{}
		err := json.Unmarshal(buf, &res)
		if err != nil {
			panic("Failed to unmarshal response: " + err.Error())
		}

		if data, ok := res.Data.(map[string]interface{}); ok {
			if tokenValue, exists := data["token"]; exists {
				token := tokenValue.(string)
			}
		}
	}
}
