package main

import (
	"encoding/json"
	"fmt"
	"gitub.com/fjqzx/xiangmu/proto"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnLineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Messager广播消息channel的foroutine，一旦有消息就发送给全部的在线User
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnLineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播信息的方法
func (s *Server) BroadCast(msg proto.Message) {
	sendMsg, _ := json.Marshal(msg)
	s.Message <- string(sendMsg)
}

func (s *Server) Handler(conn net.Conn) {
	//fmt.Printf("连接建立成功, addr:%s\n", conn.RemoteAddr().String())

	user := NewUser(conn, s)

	user.Online()

	//监听用户是否活跃的Channel
	isLive := make(chan bool)

	//接收当前用户上线消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if err != nil || n == 0 {
				user.Offline()
				return
			}

			//提取用户的消息（去除‘\n')
			msgBytes := buf[:n]
			fmt.Printf("接收到消息：%s\n", string(msgBytes))

			//用户针对msg进行消息处理
			user.DoMessage(msgBytes)

			//用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		//	当前用户是活跃的，应该重置定时器
		//不做任何事，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 3000):
			//	已经超时

			//将当前的User强制关闭
			// user.SendMsg("You got kicked")

			//销毁用的资源
			conn.Close()

			//	退出当前Handler
			return
		}
	}

	//当前handler阻塞
	select {}
}

func (s *Server) Start() {
	// 监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	defer listener.Close()

	fmt.Printf("服务开始监听，IP：%s，端口号：%d\n", s.Ip, s.Port)

	// 启动监听Message的goroutine
	go s.ListenMessager()

	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
		}

		go s.Handler(conn)
	}
}
