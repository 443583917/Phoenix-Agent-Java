package errcode

type ErrCode struct {
	Code int
	Msg  string
}

var (
	Success       = ErrCode{0, "success"}
	Unauthorized  = ErrCode{401, "未认证"}
	Forbidden     = ErrCode{403, "无权限"}
	NotFound      = ErrCode{404, "资源不存在"}
	InternalError = ErrCode{500, "服务器内部错误"}

	InvalidParams   = ErrCode{1001, "参数校验失败"}
	TooManyRequests = ErrCode{1002, "请求过于频繁"}

	AgentOffline     = ErrCode{2001, "智能体已下线"}
	AgentNotFound    = ErrCode{2002, "智能体不存在"}
	SessionExpired   = ErrCode{2003, "会话已过期"}

	DatasourceError   = ErrCode{3001, "数据源连接失败"}
	DatasourceNotFound = ErrCode{3002, "数据源不存在"}
	SQLError           = ErrCode{3003, "SQL执行失败"}

	ModelError     = ErrCode{4001, "模型调用失败"}
	EmbeddingError = ErrCode{4002, "向量化失败"}

	MilvusError  = ErrCode{5001, "向量检索失败"}
	QueueError   = ErrCode{5002, "消息队列异常"}
	CacheError   = ErrCode{5003, "缓存服务异常"}
)
