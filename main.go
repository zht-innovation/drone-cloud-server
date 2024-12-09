package main

import (
	"zhtcloud/gateway"
	"zhtcloud/models"
)

func main() {
	models.DatabaseSetup()
	gateway.ServerSetup()

	defer models.Close()
}
