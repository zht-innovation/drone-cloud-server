package response

const (
	SUCCESS        = 200
	INVALID_PARAMS = 400
	SERVER_ERROR   = 500

	ERROR_MAC_FORMAT = 600 + iota
	EXCEED_RATE_LIMIE
	INVALID_DEVICE
	INVALID_TOKEN
)
