package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		panic("请输入IP和端口号，例如：nc.exe 127.0.0.1 8080")
	}

	// 连接服务
	ip, port := args[1], args[2]
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic("连接服务失败：" + err.Error())
	}

	reader(conn)
	writer(conn)
}

// reader 读取服务端发送的消息，并打印
func reader(conn net.Conn) {
	msg := make([]byte, 1024)
	for {
		n, err := conn.Read(msg)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("服务端主动断开连接，开始退出...")
				return
			}
			panic("读取消息错误: " + err.Error())
		}
		fmt.Printf("%s\n", string(msg[:n]))
	}
}

func writer(conn net.Conn) {

}
