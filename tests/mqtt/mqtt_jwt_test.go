package mqtt_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/crypto/ssh"
)

func _setupClient(pwd string) bool {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("zht-mqtt_client")
	opts.SetUsername(pwd)
	opts.SetPassword("")
	cli := mqtt.NewClient(opts)

	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return true
}

type request struct {
	Secret string `json:"secret"`
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
			var req request

			secretBase := "zhtaero|7777777777|FF-FF-FF-FF-FF-FF"
			// 读取公钥文件
			publicKeyPath := "/home/zht/.ssh/id_rsa.pub"
			publicKeyBytes, err := os.ReadFile(publicKeyPath)
			if err != nil {
				t.Fatalf("Failed to read public key file: %v", err)
			}

			// 解析SSH公钥格式
			publicKeyStr := strings.TrimSpace(string(publicKeyBytes))
			parts := strings.Fields(publicKeyStr)
			if len(parts) < 2 {
				t.Fatal("Invalid SSH public key format")
			}

			// 解码base64编码的公钥数据
			keyData, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				t.Fatalf("Failed to decode public key: %v", err)
			}

			// 解析SSH公钥
			publicKey, err := x509.ParsePKCS1PublicKey(keyData[len("ssh-rsa")+4:])
			if err != nil {
				// 如果PKCS1失败，尝试使用ssh包解析
				t.Logf("PKCS1 parse failed, trying alternative method: %v", err)

				// 使用golang.org/x/crypto/ssh包解析SSH公钥
				sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBytes)
				if err != nil {
					t.Fatalf("Failed to parse SSH public key: %v", err)
				}

				cryptoPublicKey := sshPublicKey.(ssh.CryptoPublicKey).CryptoPublicKey()
				var ok bool
				publicKey, ok = cryptoPublicKey.(*rsa.PublicKey)
				if !ok {
					t.Fatal("Public key is not RSA")
				}
			}

			// 使用公钥加密secret
			encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(secretBase))
			if err != nil {
				t.Fatalf("Failed to encrypt secret: %v", err)
			}

			// Base64编码加密后的数据
			encryptedSecret := base64.StdEncoding.EncodeToString(encryptedData)
			t.Logf("Encrypted secret: %s", encryptedSecret)

			req.Secret = encryptedSecret
			reqbuf, _ := json.Marshal(&req)
			rsp, _ := http.Post("http://localhost:32223/auth", "application/json", bytes.NewBuffer(reqbuf))
			if rsp.StatusCode == 200 {
				buf, _ := io.ReadAll(rsp.Body)
				res := response{}
				err := json.Unmarshal(buf, &res)
				if err != nil {
					t.Error("Failed to unmarshal response:", err)
				}

				t.Log("Response msg:", res.Msg)

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
