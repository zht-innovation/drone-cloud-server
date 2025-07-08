package response

var CodeToMsgMap = map[int]string{
	SUCCESS:        "ok",
	INVALID_PARAMS: "请求参数不正确",
	SERVER_ERROR:   "fail",

	ERROR_MAC_FORMAT:  "mac地址格式不正确",
	EXCEED_RATE_LIMIE: "请求过于频繁，请稍后再试",
}
