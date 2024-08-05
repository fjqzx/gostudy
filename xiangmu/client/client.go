package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gitub.com/fjqzx/xiangmu/proto"
	"net"
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
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是:8888)")
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("4.查询在线用户")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag <= 4 && flag >= 0 {
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
		case 2:
			//fmt.Println("2.私聊模式")
			client.PrivateChat()
		case 3:
			//fmt.Println("3.修改用户名")
			client.UpdateName()
		case 4:
			client.SelectUsers()
		}
	}
}

// 处理server回应的消息，直接显示标准输出即可
func (client *Client) DealResponse() {
	//	一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	// io.Copy(os.Stdout, client.conn)

	reader := bufio.NewReader(client.conn)
	buf := make([]byte, 1024)
	for {
		// 读取服务端发送的消息
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println("conn read error: ", err)
			return
		}

		// 获取到的消息
		msgBytes := buf[:n]

		// 将消息解析成结构体
		var msg proto.Message
		_ = json.Unmarshal(msgBytes, &msg)

		// 根据消息类型进行处理消息
		switch msg.Type {
		case proto.MsgTypeOnline: // 处理上线消息
			var online proto.Online
			_ = json.Unmarshal(msg.Data, &online)
			fmt.Printf(">> 用户上线，用户名：%s | Addr:%s\n", online.Name, online.Addr)

		case proto.MsgTypeOffline: // 处理下线消息
			var offline proto.Offline
			_ = json.Unmarshal(msg.Data, &offline)
			fmt.Printf(">> 用户下线，用户名：%s | Addr:%s\n", offline.Name, offline.Addr)

		case proto.MsgTypeWho: // 处理查询用户列表消息
			var who proto.Who
			_ = json.Unmarshal(msg.Data, &who)
			fmt.Println(">> 在线用户列表")
			for _, online := range who.Onlines {
				fmt.Printf("\t用户名：%s | Addr:%s\n", online.Name, online.Addr)
			}

		case proto.MsgTypePrivate:
			var private proto.Private
			_ = json.Unmarshal(msg.Data, &private)
			fmt.Printf("用户名：%s : %s\n", private.Miname, private.Information)

		case proto.MsgTypeGroup:
			var group proto.Group
			_ = json.Unmarshal(msg.Data, &group)
			fmt.Printf("用户名：%s : %s\n", group.Miname, group.Information)

		}
	}

}

func (client *Client) SendMsg(msg proto.Message) {
	msgBytes, _ := json.Marshal(msg)
	_, _ = client.conn.Write(msgBytes)
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)

	rename := proto.Rename{Name: client.Name}
	renameBytes, _ := json.Marshal(rename)
	client.SendMsg(proto.Message{
		Type: proto.MsgTypeRename,
		Data: renameBytes,
	})

	//sendMsg := "rename|" + client.Name + "\n"
	//_, err := client.conn.Write([]byte(sendMsg))
	//if err != nil {
	//	fmt.Println("conn.Write err:", err)
	//	return false
	//}
	return true
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sdf := proto.Group{Miname: client.Name, Information: chatMsg}
			sdaBytes, _ := json.Marshal(sdf)
			Message := proto.Message{Type: proto.MsgTypeGroup, Data: sdaBytes}
			senMsg, _ := json.Marshal(Message)
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
	// 查询在线用户不需要携带任何消息，所以data为空
	client.SendMsg(proto.Message{Type: proto.MsgTypeWho})
}

// 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chaMsg string
	fmt.Println("名字：", client.Name)

	client.SelectUsers()
	fmt.Println("输入聊天对象，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("输入消息，exit退出")
		fmt.Scanln(&chaMsg)

		for remoteName != "exit" {
			if len(chaMsg) != 0 {

				//online := proto.Online{Name: u.Name, Addr: u.Addr}
				//onlineMsg, _ := json.Marshal(online)

				//发消息给服务器
				//消息不为空则发送
				//senMsg := proto.Private{Name: remoteName,Information: chaMsg}
				//senMsg := "to|" + remoteName + "|" + chaMsg + "\n\n"
				//_, err := client.conn.Write([]byte(senMsg))
				//Message := proto.Private{Name: remoteName,Information: chaMsg}\
				sdf := proto.Private{Miname: client.Name, Name: remoteName, Information: chaMsg}
				sdaBytes, _ := json.Marshal(sdf)
				Message := proto.Message{Type: proto.MsgTypePrivate, Data: sdaBytes}
				senMsg, _ := json.Marshal(Message)
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

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>No!")
		return
	}

	go client.DealResponse()

	fmt.Println(">>>>>>>>>>Ok!")

	//	启动客户端
	client.Run()
}
