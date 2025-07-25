package gateway

import (
	"fmt"
	"net/http"
	S "zhtcloud/gateway/shared"
	rsp "zhtcloud/pkg/response"
	"zhtcloud/utils"
)

// sendCoordinates handles the HTTP request to send waypoint coordinates to drones
func sendCoordinates(w http.ResponseWriter, r *http.Request) {
	rs := S.Result{}

	defer S.HandleResBodyEncode(w, &rs)

	if r.Method == http.MethodPost {
		var req S.Coordinates
		if needReturn := S.HandleReqBodyDecode(r.Body, &req, &rs); needReturn {
			return
		}

		coors := req.Coords
		redisClient := utils.GetRedisConn(S.Ctx)
		redisClient.PubChannel(S.Ctx, S.COORS, fmt.Sprintf("%v", coors))

		rs.Code = rsp.SUCCESS
		rs.Msg = rsp.CodeToMsgMap[rsp.SUCCESS]
	} else {
		S.HandleErrorReqMethod(&rs)
	}
}
