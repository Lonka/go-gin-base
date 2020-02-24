package e

var MsgFlags = map[int]string{
	SUCCESS:             "ok",
	ERROR:               "fail",
	INVALID_PARAMS:      "parameter invalid",
	ERROR_EXIST_TAG:     "tag has existed",
	ERROR_NOT_EXIST_TAG: "tag has not existed",

	ERROR_AUTH_CHECK_TOKEN_FAIL:     "check token fail",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "token timeout",
	ERROR_AUTH_TOKEN:                "generate token fail",
	ERROR_AUTH:                      "auth error",
	ERROR_AUTH_NOT_FOUND:            "token not found",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "save image fail",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "check image fail",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "image format fail",

	ERROR_EXPORT_FAIL: "export fail",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
