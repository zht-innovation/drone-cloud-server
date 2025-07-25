package emqx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	S "zhtcloud/gateway/shared"
	rsp "zhtcloud/pkg/response"
)

type topic struct {
	Node    string `json:"node"`
	Session string `json:"session"`
	Topic   string `json:"topic"`
}

type topicResponse struct {
	Meta   Meta    `json:"meta"`
	Topics []topic `json:"data"`
}

func GetEmqxTopics(token string) (topicResponse, error) {
	baseURL := S.MQTT_BROKER + S.API_PREFIX + "/topics"
	params := url.Values{}
	params.Add("node", "emqx@110.42.101.86")
	params.Add("page", "1")
	params.Add("limit", "50")

	req, err := http.NewRequest(http.MethodGet, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return topicResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return topicResponse{}, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return topicResponse{}, errors.New("failed to fetch topics, status code: " + rsp.Status)
	}

	body, _ := io.ReadAll(rsp.Body)
	var topicRsp topicResponse
	if err := json.Unmarshal(body, &topicRsp); err != nil {
		return topicResponse{}, err
	}

	return topicRsp, nil
}

func GetTopicsList(w http.ResponseWriter, r *http.Request) {
	rs := S.Result{}
	defer S.HandleResBodyEncode(w, &rs)

	if r.Method == http.MethodGet {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]
		topics, err := GetEmqxTopics(token)
		if err != nil {
			rs.Code = rsp.SERVER_ERROR
			rs.Msg = rsp.CodeToMsgMap[rsp.SERVER_ERROR]
			return
		}

		rs.Code = rsp.SUCCESS
		rs.Msg = rsp.CodeToMsgMap[rsp.SUCCESS]
		iData := interface{}(
			map[string]topicResponse{
				"topics": topics,
			},
		)
		rs.Data = &iData
	} else {
		S.HandleErrorReqMethod(&rs)
	}
}
