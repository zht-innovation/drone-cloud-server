package mqtt_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func _setupClient(pwd string) bool {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(pwd)
	opts.SetPassword("")
	cli := mqtt.NewClient(opts)

	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return true
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// use 'go test -v' to run this test
func TestMqttJwt(t *testing.T) {
	// catch panic in '_setupClient' if it fails to connect
	defer func() {
		if r := recover(); r != nil {
			t.Log("Failed to setup MQTT client for No JWT")
			rsp, _ := http.Get("http://localhost:32223/auth?mac=00-11-22-33-44-55")
			if rsp.StatusCode == 200 {
				buf, _ := io.ReadAll(rsp.Body)
				res := response{}
				err := json.Unmarshal(buf, &res)
				if err != nil {
					t.Error("Failed to unmarshal response:", err)
				}

				if data, ok := res.Data.(map[string]interface{}); ok {
					if tokenValue, exists := data["token"]; exists {
						token := tokenValue.(string)
						t.Log("Received JWT token:", token)
						ok := _setupClient(token)
						if !ok {
							t.Error("Failed to setup MQTT client with JWT")
						} else {
							t.Log("MQTT client setup successful with JWT")
						}
					}
				}
			}
		} else {
			t.Log("MQTT client setup successful for No JWT")
		}
	}()

	_setupClient("")
}
