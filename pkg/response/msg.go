package response

var CodeToMsgMap = map[int]string{
	SUCCESS:        "ok",
	INVALID_PARAMS: "请求参数不正确",
	SERVER_ERROR:   "fail",
}
