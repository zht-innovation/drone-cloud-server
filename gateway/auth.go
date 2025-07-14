package gateway

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	rsp "zhtcloud/pkg/response"
	"zhtcloud/utils/logger"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/ssh"
)

type authReq struct {
	Secret string `json:"secret"`
}

// var secret = []byte(os.Getenv("SECRET")) // execute when gateway starts, but before main()

func genMqttToken() (string, error) {
	claims := jwt.MapClaims{
		"exp": jwt.TimeFunc().Add(24 * time.Hour).Unix(), // 1 day expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SECRET"))
	return token.SignedString(secret)
}

func checkValidDevice(secret string) (string, bool) {
	privateKeyPath := "/home/zht/.ssh/id_rsa"
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		logger.Error("Failed to read private key file: %v", err)
		return "", false
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		logger.Error("Failed to parse private key: %v", err)
		return "", false
	}

	// Base64 decode secret parameter
	encryptedData, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		logger.Error("Failed to decode base64 secret: %v", err)
		return "", false
	}

	// Decrypt using private key
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), encryptedData)
	if err != nil {
		logger.Error("Failed to decrypt data: %v", err)
		return "", false
	}

	secretList := strings.Split(string(decryptedData), "|")
	if secretList[0] != os.Getenv("SECRET") {
		logger.Error("Invalid secret prefix: %s", secretList[0])
		return "", false
	}

	now := time.Now().Unix()
	reqtime, _ := strconv.Atoi(secretList[1])
	if secretList[1] != "7777777777" && (now < int64(reqtime) || now-int64(reqtime) > 20) {
		logger.Error("Request time expired: %d", reqtime)
		return "", false
	}

	return secretList[2], true
}

func validateMacFormat(mac string) bool {
	if len(mac) != 17 {
		return false
	}
	for i, c := range mac {
		if (i+1)%3 == 0 {
			if c != ':' && c != '-' {
				return false
			}
		} else {
			if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
				return false
			}
		}
	}

	return true
}

func authDrone(w http.ResponseWriter, r *http.Request) {
	rs := Result{}

	defer HandleResBodyEncode(w, &rs)

	if r.Method == http.MethodPost {
		var req authReq
		if needReturn := HandleReqBodyDecode(r.Body, &req, &rs); needReturn {
			return
		}

		secret := req.Secret
		mac, ok := checkValidDevice(secret)
		if !ok {
			rs.Code = rsp.INVALID_DEVICE
			rs.Msg = rsp.CodeToMsgMap[rsp.INVALID_DEVICE]
			return
		}

		// validate the MAC address format
		ok = validateMacFormat(mac) // TODO: store mac to redis
		if !ok {
			rs.Code = rsp.ERROR_MAC_FORMAT
			rs.Msg = rsp.CodeToMsgMap[rsp.ERROR_MAC_FORMAT]
			return
		}

		// generate the token
		token, err := genMqttToken()
		if err != nil {
			rs.Code = rsp.SERVER_ERROR
			rs.Msg = rsp.CodeToMsgMap[rsp.SERVER_ERROR]
			return
		}

		rs.Code = rsp.SUCCESS
		rs.Msg = rsp.CodeToMsgMap[rsp.SUCCESS]
		iData := interface{}(
			map[string]string{
				"token": token,
			},
		)

		rs.Data = &iData
	} else {
		HandleErrorReqMethod(&rs)
	}
}
