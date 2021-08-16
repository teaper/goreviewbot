package msgc

/*
用以存储 update 中出现的验证码和用户的绑定，方便后期删除和踢出用户
 */
type AuthUser struct {
	UserID int
	UserName string
}

type CodeMessage struct {
	MessageID int  //验证码消息的 id
	Codes [4] int32  //验证码
	AuthUser AuthUser //验证用户
	Enabled bool //消息状态，存在（true）还是删除（false）
}

