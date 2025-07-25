package gateway

import (
	"net/http"
	"zhtcloud/gateway/emqx"
	mw "zhtcloud/middleware"
	"zhtcloud/utils/logger"
)

func ServerSetup() {
	http.HandleFunc("/drone", droneDataHandler)
	http.HandleFunc("/frontend", frontendTransfer)
	http.HandleFunc("/coords", sendCoordinates)
	http.HandleFunc("/auth", mw.CORSMiddleWare(authDrone))

	http.HandleFunc("/emqx/login", mw.CORSMiddleWare(emqx.AdminLogin))

	err := http.ListenAndServe("0.0.0.0:32223", nil)
	if err != nil {
		logger.Fatal("http listen and serve %v", err)
	}
}
