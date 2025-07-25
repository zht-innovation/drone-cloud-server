package response

var CodeToMsgMap = map[int]string{
	SUCCESS:        "ok",
	INVALID_PARAMS: "请求参数不正确",
	SERVER_ERROR:   "服务端错误",

	ERROR_MAC_FORMAT:  "mac地址格式不正确",
	EXCEED_RATE_LIMIE: "请求过于频繁，请稍后再试",
	INVALID_DEVICE:    "设备校验失败",
	INVALID_TOKEN:     "无效的令牌",
}
