package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	//IP
	ServerIP string
	//端口
	ServerPort int
	//名字
	Name string
	//特别字句柄
	conn net.Conn
	//当前client的模式
	flag int
}

func NewClient(serverIP string, serverPort int) *Client {
	//	创建客户端对象
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       999,
	}
	//	链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn
	//	返回对象
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器地址(默认是:127.0.0.1)")
	flag.IntVar(&serverPort, "Port", 8888, "设置服务器端口(默认是:8888)")
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag <= 3 && flag >= 0 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的数字！")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			//fmt.Println("1.公聊模式")
			client.PublicChat()
			break
		case 2:
			//fmt.Println("2.私聊模式")
			client.PrivateChat()
			break
		case 3:
			//fmt.Println("3.修改用户名")
			client.UpdateName()
			break
		}
	}
}

// 处理server回应的消息，直接显示标准输出即可
func (client *Client) DealResponse() {
	//	一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			//发消息给服务器

			//消息不为空则发送
			senMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(senMsg))
			if err != nil {
				fmt.Println("conn.Write err", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

// 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return
	}

}

// 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chaMsg string

	client.SelectUsers()
	fmt.Println("输入聊天对象，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("输入消息，exit退出")
		fmt.Scanln(&chaMsg)

		for remoteName != "exit" {
			if len(chaMsg) != 0 {
				//发消息给服务器

				//消息不为空则发送
				senMsg := "to|" + remoteName + "|" + chaMsg + "\n\n"
				_, err := client.conn.Write([]byte(senMsg))
				if err != nil {
					fmt.Println("conn.Write err", err)
					break
				}
			}
			chaMsg = ""
			fmt.Println("输入消息，exit退出")
			fmt.Scanln(&chaMsg)
		}
	}
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>>No!")
		return
	}

	go client.DealResponse()

	fmt.Println(">>>>>>>>>>Ok!")

	//	启动客户端
	client.Run()
}
