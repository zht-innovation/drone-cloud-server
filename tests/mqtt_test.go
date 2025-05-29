package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var topic = "drone/data"

var sigChan = make(chan struct{})
var msgPayloadRecv string
var msgPayloadSend string
var msgPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(msg.Payload())
	msgPayloadRecv = string(msg.Payload())
}
var connHandler mqtt.OnConnectHandler = func(client mqtt.Client) {}
var connLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {}

func setupClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.OnConnect = connHandler
	opts.OnConnectionLost = connLostHandler
	cli := mqtt.NewClient(opts)

	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return cli
}

func subscriber(cli mqtt.Client) {
	token := cli.Subscribe(topic, 1, msgPubHandler)
	token.Wait()

	// 保持连接，直到收到第一条测试信息
	<-sigChan
}

func publisher(cli mqtt.Client) {
	lat := 30.0 + rand.Float64()*0.1
	lng := 120.0 + rand.Float64()*0.1
	alt := 100.0 + rand.Float64()*10.0

	payload := fmt.Sprintf(`{
		"location": {
			"lat": %.6f,
			"lng": %.6f,
			"alt": %.2f
		}
	}`, lat, lng, alt)

	msgPayloadSend = payload

	token := cli.Publish(topic, 0, false, payload)
	token.Wait()

	sigChan <- struct{}{}
}

func TestMqtt(t *testing.T) {
	cli := setupClient() // 订阅方和发布方使用同一个client

	go subscriber(cli)
	time.Sleep(time.Second * 2)

	publisher(cli)
	time.Sleep(time.Second * 2)

	if msgPayloadRecv != msgPayloadSend {
		t.Errorf("Received payload: %s does not match sent payload: %s", msgPayloadRecv, msgPayloadSend)
	} else {
		t.Logf("Received payload: %s matches sent payload: %s", msgPayloadRecv, msgPayloadSend)
	}

	defer cli.Unsubscribe(topic)
	defer cli.Disconnect(250)
}
