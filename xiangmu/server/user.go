package main

import (
	"encoding/json"
	"fmt"
	"gitub.com/fjqzx/xiangmu/proto"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user chanel消息的goroutine
	go user.ListenMessage()

	return &user
}

// 用户上线服务
func (u *User) Online() {
	//用户上线将用户添加到OnLineMap中
	u.server.mapLock.Lock()
	u.server.OnLineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播当前用户上线信息
	online := proto.Online{Name: u.Name, Addr: u.Addr}
	onlineMsg, _ := json.Marshal(online)
	u.server.BroadCast(proto.Message{
		Type: proto.MsgTypeOnline,
		Data: onlineMsg,
	})
}

// 用户下线服务
func (u *User) Offline() {
	//用户上线将用户从OnLineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnLineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播当前用户下线信息
	offline := proto.Offline{Name: u.Name, Addr: u.Addr}
	offlineMsg, _ := json.Marshal(offline)
	u.server.BroadCast(proto.Message{
		Type: proto.MsgTypeOffline,
		Data: offlineMsg,
	})
}

// 给当前User对应的客户端发消息
func (u *User) SendMsg(msg proto.Message) {
	// 将结构体转换成切片
	msgBytes, _ := json.Marshal(msg)
	// 发送消息
	_, err := u.conn.Write(msgBytes)
	if err != nil {
		fmt.Printf("消息发送失败，用户名：%s | 错误信息：%v\n", u.Name, err)
	}
}

// 用户处理消息的业务
func (u *User) DoMessage(msgBytes []byte) {
	var msg proto.Message
	_ = json.Unmarshal(msgBytes, &msg)

	switch msg.Type {
	case proto.MsgTypeWho:
		// 查询所有在线用户
		who := proto.Who{}
		u.server.mapLock.Lock()
		for _, user := range u.server.OnLineMap {
			who.Onlines = append(who.Onlines, proto.Online{Name: user.Name, Addr: user.Addr})
		}
		u.server.mapLock.Unlock()
		// 发送所有在线用户
		whoBytes, _ := json.Marshal(who)
		u.SendMsg(proto.Message{
			Type: proto.MsgTypeWho,
			Data: whoBytes,
		})

	case proto.MsgTypeRename:
		// 获取到重命名消息
		var rename proto.Rename
		_ = json.Unmarshal(msg.Data, &rename)

		// 处理重命名操作：判断name是否存在
		_, ok := u.server.OnLineMap[rename.Name]
		if ok {
			u.SendMsg(proto.Message{Type: proto.MsgTypeError, Data: []byte("当前用户名已经被使用")})
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnLineMap, u.Name)
			u.server.OnLineMap[rename.Name] = u
			u.server.mapLock.Unlock()
			u.Name = rename.Name
			u.SendMsg(proto.Message{Type: proto.MsgTypeOk, Data: []byte(fmt.Sprintf("修改用户名 %s 成功", u.Name))})
		}
	case proto.MsgTypePrivate:
		var privMsg proto.Private

		_ = json.Unmarshal(msg.Data, &privMsg)
		//remoteUser, ok := u.server.OnLineMap[remoteName]
		remoteUser, ok := u.server.OnLineMap[privMsg.Name]
		if !ok {
			u.SendMsg(proto.Message{Type: proto.MsgTypeError, Data: []byte("name No!")})
			return
		}
		if privMsg.Information == "" {
			u.SendMsg(proto.Message{Type: proto.MsgTypeError, Data: []byte("message No!")})
			return
		}
		remoteUser.SendMsg(msg)

	case proto.MsgTypeGroup:
		u.server.BroadCast(msg)

	case proto.MsgTypeTransfer:
		var transMsg proto.Transfer

		_ = json.Unmarshal(msg.Data, &transMsg)
		//remoteUser, ok := u.server.OnLineMap[remoteName]
		remoteUser, ok := u.server.OnLineMap[transMsg.Name]
		fmt.Println("找到文件接收用户：", remoteUser.Name)
		if !ok {
			u.SendMsg(proto.Message{Type: proto.MsgTypeError, Data: []byte("name No!")})
			return
		}
		remoteUser.SendMsg(msg)
	}
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		//fmt.Printf("广播消息，addr:%s msg:%s\n", u.Addr, msg)
		u.conn.Write([]byte(msg))
	}
}
