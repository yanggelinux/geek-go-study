package ce

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_RESOURCE_EXIST     = 10001
	ERROR_RESOURCE_NOT_EXIST = 10002

	ERROR_AUTH_CHECK_TOKEN_FAILED  = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN_FAILED        = 20003
	ERROR_AUTH_TOKEN               = 20004
	ERROR_AUTH_CHECK_FAILED        = 20005

	ERROR_DB_TABLE_QUERY_FAILED     = 3001
	ERROR_DB_TABLE_CREATE_FAILED    = 3002
	ERROR_DB_TABLE_UPDATE_FAILED    = 3003
	ERROR_DB_TABLE_DELETE_FAILED    = 3004
	ERROR_DB_TABLE_RECORD_NOT_EXIST = 3005

	ERROR_USER_NOT_EXIST = 4002
	ERROR_PASSWORD       = 4003
)

var MsgFlags = map[int]string{
	SUCCESS:        "请求成功",
	ERROR:          "请求错误",
	INVALID_PARAMS: "请求参数错误",

	ERROR_RESOURCE_EXIST:     "资源已存在",
	ERROR_RESOURCE_NOT_EXIST: "资源不存在",

	ERROR_AUTH_CHECK_TOKEN_FAILED:  "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN_FAILED:        "请求Token缺失",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH_CHECK_FAILED:        "鉴权失败",

	ERROR_DB_TABLE_QUERY_FAILED:     "数据库表查询失败",
	ERROR_DB_TABLE_CREATE_FAILED:    "数据库表记录创建失败",
	ERROR_DB_TABLE_UPDATE_FAILED:    "数据库表记录更新失败",
	ERROR_DB_TABLE_DELETE_FAILED:    "数据库表记录删除失败",
	ERROR_DB_TABLE_RECORD_NOT_EXIST: "数据库表中记录不存在",

	ERROR_USER_NOT_EXIST: "用户不存在",
	ERROR_PASSWORD:       "密码错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

type Data struct {
}

type Respone struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   *Data  `json:"data"`
}
