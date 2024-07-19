package main

import (
	"net"
	"strings"
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
func (this *User) Online() {
	//用户上线将用户添加到OnLineMap中
	this.server.mapLock.Lock()
	this.server.OnLineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播当前用户上线信息
	this.server.BroadCast(this, "hello")
}

// 用户下线服务
func (this *User) Offline() {
	//用户上线将用户从OnLineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnLineMap, this.Name)
	this.server.mapLock.Unlock()
	//广播当前用户下线信息
	this.server.BroadCast(this, "baibai")
}

// 给当前User对应的客户端发消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//	查询当前在线用户都有那些
		this.server.mapLock.Lock()
		for _, user := range this.server.OnLineMap {
			onlineMsg := "{" + user.Addr + "}" + user.Name + ":" + "在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//信息格式：rename|张三
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := this.server.OnLineMap[newName]
		if ok {
			this.SendMsg("当前用户名已经被使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnLineMap, this.Name)
			this.server.OnLineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("Your name changed successfully!" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("to|name|message")
			return
		}
		//根据对方的用户名，获取他的User对象
		remoteUser, ok := this.server.OnLineMap[remoteName]
		if !ok {
			this.SendMsg("name NO!!")
			return
		}
		//获取信息内容，通过他的User对象将内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("message No!")
			return
		}
		remoteUser.SendMsg(this.Name + ":" + content)
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		//fmt.Printf("广播消息，addr:%s msg:%s\n", this.Addr, msg)
		this.conn.Write([]byte(msg + "\n"))
	}
}
