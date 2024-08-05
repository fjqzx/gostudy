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
	_, _ = u.conn.Write(msgBytes)
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
	}

	//// ！！！ 开始改造第一个协议 who
	//if msg == "who" {
	//	//	查询当前在线用户都有那些
	//	u.server.mapLock.Lock()
	//	for _, user := range u.server.OnLineMap {
	//		onlineMsg := "{" + user.Addr + "}" + user.Name + ":" + "在线...\n"
	//		u.SendMsg(onlineMsg)
	//	}
	//	u.server.mapLock.Unlock()
	//} else if len(msg) > 7 && msg[:7] == "rename|" {
	//	//信息格式：rename|张三
	//	newName := strings.Split(msg, "|")[1]
	//	//判断name是否存在
	//	_, ok := u.server.OnLineMap[newName]
	//	if ok {
	//		u.SendMsg("当前用户名已经被使用\n")
	//	} else {
	//		u.server.mapLock.Lock()
	//		delete(u.server.OnLineMap, u.Name)
	//		u.server.OnLineMap[newName] = u
	//		u.server.mapLock.Unlock()
	//		u.Name = newName
	//		u.SendMsg("Your name changed successfully!" + u.Name + "\n")
	//	}
	//} else if len(msg) > 4 && msg[:3] == "to|" {
	//	//获取对方的用户名
	//	remoteName := strings.Split(msg, "|")[1]
	//	if remoteName == "" {
	//		u.SendMsg("to|name|message")
	//		return
	//	}
	//	//根据对方的用户名，获取他的User对象
	//	remoteUser, ok := u.server.OnLineMap[remoteName]
	//	if !ok {
	//		u.SendMsg("name NO!!")
	//		return
	//	}
	//	//获取信息内容，通过他的User对象将内容发送过去
	//	content := strings.Split(msg, "|")[2]
	//	if content == "" {
	//		u.SendMsg("message No!")
	//		return
	//	}
	//	remoteUser.SendMsg(u.Name + ":" + content)
	//} else {
	//	u.server.BroadCast(u, msg)
	//}
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		//fmt.Printf("广播消息，addr:%s msg:%s\n", u.Addr, msg)
		u.conn.Write([]byte(msg))
	}
}
