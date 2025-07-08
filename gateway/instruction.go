package gateway

import (
	"fmt"
	"net/http"
	rsp "zhtcloud/pkg/response"
	"zhtcloud/utils"
)

// sendCoordinates handles the HTTP request to send waypoint coordinates to drones
func sendCoordinates(w http.ResponseWriter, r *http.Request) {
	rs := Result{}

	defer HandleResBodyEncode(w, &rs)

	if r.Method == http.MethodPost {
		var req Coordinates
		if needReturn := HandleReqBodyDecode(r.Body, &req, &rs); needReturn {
			return
		}

		coors := req.Coords
		redisClient := utils.GetRedisConn(ctx)
		redisClient.PubChannel(ctx, COORS, fmt.Sprintf("%v", coors))

		rs.Code = rsp.SUCCESS
		rs.Msg = rsp.CodeToMsgMap[rsp.SUCCESS]
	} else {
		HandleErrorReqMethod(&rs)
	}
}
