package emqx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	S "zhtcloud/gateway/shared"
	rsp "zhtcloud/pkg/response"

	"github.com/golang-jwt/jwt/v4"
)

type licence struct {
	Edition string `json:"edition"`
}

type adminLoginResponse struct {
	Version string  `json:"version"`
	Role    string  `json:"role"`
	Token   string  `json:"token"`
	License licence `json:"license"`
}

func emqxLogin() (string, error) {
	username := os.Getenv("EMQX_USERNAME")
	password := os.Getenv("EMQX_PASSWORD")

	url := fmt.Sprintf("%s%s/login", S.MQTT_BROKER, S.API_PREFIX)
	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return "", errors.New("failed to login to EMQX, status code: " + rsp.Status)
	}

	rspBody, _ := io.ReadAll(rsp.Body)
	var loginRsp adminLoginResponse
	if err := json.Unmarshal(rspBody, &loginRsp); err != nil {
		return "", err
	}

	return loginRsp.Token, nil
}

func validateToken(tokenString string) bool {
	secret := []byte(os.Getenv("SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	} else {
		return false
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	rs := S.Result{}

	defer S.HandleResBodyEncode(w, &rs)

	if r.Method == http.MethodPost {
		token := r.Header.Get("Token")
		if token == "" || !validateToken(token) {
			rs.Code = rsp.INVALID_TOKEN
			rs.Msg = rsp.CodeToMsgMap[rsp.INVALID_TOKEN]
			return
		}

		adminToken, err := emqxLogin()
		if err != nil {
			rs.Code = rsp.SERVER_ERROR
			rs.Msg = rsp.CodeToMsgMap[rsp.SERVER_ERROR]
			return
		}

		rs.Code = rsp.SUCCESS
		rs.Msg = rsp.CodeToMsgMap[rsp.SUCCESS]
		iData := interface{}(
			map[string]string{
				"adminToken": adminToken,
			},
		)
		rs.Data = &iData
	} else {
		S.HandleErrorReqMethod(&rs)
	}
}
