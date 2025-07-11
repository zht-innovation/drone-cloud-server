package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	r "math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/crypto/ssh"
)

type request struct {
	Secret string `json:"secret"`
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func _setupClient(pwd string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ws://110.42.101.86:8083/mqtt")
	opts.SetClientID("zht-websocket1")
	opts.SetUsername(pwd)
	opts.SetPassword("")
	cli := mqtt.NewClient(opts)

	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return cli
}

func encrypt(base string) (string, error) {
	// 读取公钥文件
	publicKeyPath := "/home/zht/.ssh/id_rsa.pub"
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	// 解析SSH公钥格式
	publicKeyStr := strings.TrimSpace(string(publicKeyBytes))
	parts := strings.Fields(publicKeyStr)
	if len(parts) < 2 {
		return "", err
	}

	// 解码base64编码的公钥数据
	keyData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// 解析SSH公钥
	publicKey, err := x509.ParsePKCS1PublicKey(keyData[len("ssh-rsa")+4:])
	if err != nil {
		// 使用golang.org/x/crypto/ssh包解析SSH公钥
		sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBytes)
		if err != nil {
			return "", err
		}

		cryptoPublicKey := sshPublicKey.(ssh.CryptoPublicKey).CryptoPublicKey()
		var ok bool
		publicKey, ok = cryptoPublicKey.(*rsa.PublicKey)
		if !ok {
			return "", err
		}
	}

	// 使用公钥加密secret
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(base))
	if err != nil {
		return "", err
	}

	// Base64编码加密后的数据
	encryptedSecret := base64.StdEncoding.EncodeToString(encryptedData)
	return encryptedSecret, nil
}

func main() {
	var req request
	req.Secret, _ = encrypt("zhtaero|7777777777|01-23-01-23-01-23")
	reqbuf, _ := json.Marshal(&req)
	rsp, _ := http.Post("http://localhost:32223/auth", "application/json", bytes.NewBuffer(reqbuf))
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
				cli := _setupClient(token)
				topic := "drone/data"
				fmt.Println("Start mocking...")
				for {
					// 生成模拟数据
					time_boot_ms := 0.0 + r.Float64()*1000
					lat := 30.0 + r.Float64()*0.1
					lon := 120.0 + r.Float64()*0.1
					alt := 50.0 + r.Float64()*10
					relative_alt := 10.0 + r.Float64()*5
					vx, vy, vz := 0, 0, 0
					hdg := 0

					pitch := -50.0 + r.Float64()*110
					yaw := -50 + r.Float64()*100

					GLOBAL_POSITION_INT := fmt.Sprintf(`{
                        "time_boot_ms": %.2f,
                        "lat": %.6f,
                        "lon": %.6f,
                        "alt": %.2f,
                        "relative_alt": %.2f,
                        "vx": %d,
                        "vy": %d,
                        "vz": %d,
                        "hdg": %d
                	}`, time_boot_ms, lat, lon, alt, relative_alt, vx, vy, vz, hdg)

					ATTITUDE := fmt.Sprintf(`{
                        "roll": 0,
                        "pitch": %.2f,
                        "yaw": %.2f,
                        "rollspeed": 0,
                        "pitchspeed": 0,
                        "yawspeed": 0
                	}`, pitch, yaw)

					SYS_STATUS := `{
                        "onboard_control_sensors_present": 0,
                        "onboard_control_sensors_enabled": 0,
                        "onboard_control_sensors_health": 0,
                        "load": 0,
                        "voltage_battery": 12.5,
                        "current_battery": 0,
                        "battery_remaining": 100,
                        "drop_rate_comm": 0,
                        "errors_comm": 0,
                        "errors_count1": 0,
                        "errors_count2": 0,
                        "errors_count3": 0,
                        "errors_count4": 0
                	}`

					payload := fmt.Sprintf(`{
                        "GLOBAL_POSITION_INT": %s,
                        "ATTITUDE": %s,
                        "SYS_STATUS": %s,
                        "MODE": 0,
                        "STATUS": 0
                	}`, GLOBAL_POSITION_INT, ATTITUDE, SYS_STATUS)

					token := cli.Publish(topic, 1, false, payload)
					token.Wait()
					// fmt.Printf("Published: %s\n", payload)

					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}
