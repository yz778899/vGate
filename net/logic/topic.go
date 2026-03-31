package logic

const (
	//客户端上线通知
	Notice_On_Line string = "online"
	//客户端下线通知
	Notice_Off_Line string = "offline"
	//session_id 变更 [一般用户登录后可以通知网关主动变更，但需要保证其唯一性，]
	Session_Id_Change string = "session_id_change"
)
