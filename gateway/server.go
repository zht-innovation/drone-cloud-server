package gateway

import (
	"log"
	"net/http"
)

func ServerSetup() {
	http.HandleFunc("/drone", droneDataHandler)

	err := http.ListenAndServe("0.0.0.0:32223", nil)
	if err != nil {
		log.Fatalf("http listen and serve %v", err)
	}
}