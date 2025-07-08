package gateway

import (
	"net/http"
	rsp "zhtcloud/pkg/response"

	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte("zhtaero")

func genMqttToken(mac string) (string, error) {
	claims := jwt.MapClaims{
		"mac": mac,
		"exp": jwt.TimeFunc().Add(1 * 60 * 60).Unix(), // 1 hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
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

	if r.Method == http.MethodGet {
		mac := r.URL.Query().Get("mac")

		// validate the MAC address format
		ok := validateMacFormat(mac)
		if !ok {
			rs.Code = rsp.ERROR_MAC_FORMAT
			rs.Msg = rsp.CodeToMsgMap[rsp.ERROR_MAC_FORMAT]
			return
		}

		// generate the token
		token, err := genMqttToken(mac)
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
