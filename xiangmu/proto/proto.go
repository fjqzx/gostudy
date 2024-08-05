package proto

const (
	MsgTypeError   = -1 // 错误消息类型
	MsgTypeOk      = 0  // 成功消息类型
	MsgTypeWho     = 1  // 查询在线用户消息类型
	MsgTypeOnline  = 2  // 上线消息类型
	MsgTypeRename  = 3  // 重命名消息类型
	MsgTypeOffline = 4  // 下线消息
	MsgTypePrivate = 5  // 私密消息类型
	MsgTypeGroup   = 6  //群聊
)

// Message 标准消息格式
type Message struct {
	Type int    `json:"type"` // 消息类型
	Data []byte `json:"data"` // 消息内容
}

// Who 所有在线用户
type Who struct {
	Onlines []Online `json:"onlines"` // 所有在线用户信息
}

// Online 上线消息
type Online struct {
	Name string `json:"name"` // 用户名
	Addr string `json:"addr"` // 地址
}

// Offline 下线消息
type Offline struct {
	Name string `json:"name"` // 用户名
	Addr string `json:"addr"` // 地址
}

// Rename 重命名消息
type Rename struct {
	Name string `json:"name"` // 用户名
}

type Private struct {
	Miname      string `json:"miname"`      //发送人的名字
	Name        string `json:"name"`        //	接收人的名字
	Information string `json:"information"` // 消息内容
}

type Group struct {
	Miname      string `json:"miname"`      //发送人的名字
	Information string `json:"information"` // 消息内容
}
