package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gitub.com/fjqzx/xiangmu/proto"
	"io"
	"net"
	"os"
	"path/filepath"
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
	client.Name = conn.LocalAddr().String()
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
	fmt.Println("5.文件传输")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag <= 5 && flag >= 0 {
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
			// 查询在线用户
			client.SelectUsers()
		case 5:
			//文件传输
			client.File()
		}
	}
}

// 处理server回应的消息，直接显示标准输出即可
func (client *Client) DealResponse() {
	//	一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	// io.Copy(os.Stdout, client.conn)

	reader := bufio.NewReader(client.conn)
	buf := make([]byte, 4096)
	for {
		// 读取服务端发送的消息
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println("conn read error: ", err)
			os.Exit(1)
		}

		// 获取到的消息
		msgBytes := buf[:n]
		fmt.Printf("消息大小：%d 内容：%s\n", n, string(msgBytes))

		// 将消息解析成结构体
		var msg proto.Message
		_ = json.Unmarshal(msgBytes, &msg)

		// 根据消息类型进行处理消息
		switch msg.Type {
		case proto.MsgTypeError: // 处理错误消息
			fmt.Printf("[ERROR] 错误消息：%s\n", string(msg.Data))

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

		case proto.MsgTypeTransfer:
			var transfer proto.Transfer

			_ = json.Unmarshal(msg.Data, &transfer)
			fmt.Printf("用户名：%s ,向您传输文件。\n", transfer.Miname)

			// 默认存放在 data 目录
			_ = os.MkdirAll("./data", 777)
			a := filepath.Join("data", transfer.Filename)

			//fmt.Print("请输入文件存放地址:")
			//var f string
			//fmt.Scanln(&f)
			////a := filepath.Join(f, transfer.Filename)
			//a := f + "\\" + transfer.Filename

			fa, err := os.Create(a)
			if err != nil {
				fmt.Println(err)
				return
			}

			l, err := fa.WriteString(string(transfer.Information))
			if err != nil {
				fmt.Println(err)
				fa.Close()
				return
			}
			fmt.Println(l, "bytes written successfully")
			err = fa.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			/// H:\2024\gostudy\kkk.go
			//for {
			//	fmt.Println("同意请输 1 ，不同意请输 2")
			//	var a int
			//	fmt.Scanln(&a)
			//	if a == 1 {
			//		fmt.Println("yse")
			//		sdf := proto.Private{Miname: transfer.Name, Name: transfer.Miname, Information: "同意您的文件传输申请！"}
			//		sdaBytes, _ := json.Marshal(sdf)
			//		Message := proto.Message{Type: proto.MsgTypePrivate, Data: sdaBytes}
			//		senMsg, _ := json.Marshal(Message)
			//		_, err := client.conn.Write([]byte(senMsg))
			//		if err != nil {
			//			fmt.Println("conn.Write err", err)
			//			//break
			//		}
			//
			//		fmt.Println("请输入文件存放地址")
			//		var f string
			//		fmt.Scanln(&f)
			//		a := f + "\\" + transfer.Filename
			//
			//		fa, err := os.Create(a)
			//		if err != nil {
			//			fmt.Println(err)
			//			return
			//		}
			//
			//		l, err := fa.WriteString(string(transfer.Information))
			//		if err != nil {
			//			fmt.Println(err)
			//			fa.Close()
			//			return
			//		}
			//		fmt.Println(l, "bytes written successfully")
			//		err = fa.Close()
			//		if err != nil {
			//			fmt.Println(err)
			//			return
			//		}
			//		return
			//	} else if a == 2 {
			//		sdf := proto.Private{Miname: transfer.Name, Name: transfer.Miname, Information: "不同意您的文件传输申请！"}
			//		sdaBytes, _ := json.Marshal(sdf)
			//		Message := proto.Message{Type: proto.MsgTypePrivate, Data: sdaBytes}
			//		senMsg, _ := json.Marshal(Message)
			//		_, err := client.conn.Write([]byte(senMsg))
			//		if err != nil {
			//			fmt.Println("conn.Write err", err)
			//			break
			//		}
			//		return
			//	} else {
			//		fmt.Println("输入合法的数字")
			//	}
			//}
		}
	}

}

func (client *Client) File() {
	client.SelectUsers()
	fmt.Println("请选择文件传输对象，exit退出")
	var b string
	fmt.Scanln(&b)

	if b != "exit" {
		fmt.Print("请输入文件路径:")
		var filename string
		fmt.Scanln(&filename)

		_file := filepath.Base(filename)
		fmt.Println("文件名：", _file)
		fil, _ := os.Open(filename)
		defer fil.Close()
		//// 使用Stat函数获取文件信息
		//fileInfo, er := os.Stat(filename)
		//if er != nil {
		//	// 如果文件不存在或发生其他错误
		//	fmt.Println("Error:", er)
		//	return
		//}
		//
		//// 使用Size()方法获取文件大小
		//fileSize := fileInfo.Size()
		//
		//fmt.Println("文件大小：", fileSize)
		buf := make([]byte, 4096)

		var bytes []byte
		for {
			count, err := fil.Read(buf)

			if err == io.EOF {
				break
			}

			currBytes := buf[:count]

			bytes = append(bytes, currBytes...)
		}

		sdf := proto.Transfer{Miname: client.Name, Name: b, Filename: _file, Information: bytes}
		sdaBytes, _ := json.Marshal(sdf)
		Message := proto.Message{Type: proto.MsgTypeTransfer, Data: sdaBytes}
		senMsg, _ := json.Marshal(Message)
		_, err := client.conn.Write(senMsg)
		// 将字节切片转为字符串 最后打印出来文件内容
		fmt.Println(string(sdf.Information))

		if err != nil {
			fmt.Println("conn.Write err", err)
			return
		}

	}

}

func (client *Client) SendMsg(msg proto.Message) {
	msgBytes, _ := json.Marshal(msg)
	_, _ = client.conn.Write(msgBytes)
}

func (client *Client) UpdateName() bool {
	fmt.Print("请输入用户名：")
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
