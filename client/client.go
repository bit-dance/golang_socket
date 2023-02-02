package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverip string, serverport int) *Client {
	client := &Client{
		ServerIp:   serverip,
		ServerPort: serverport,
		flag:       214748,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverip, serverport))
	if err != nil {
		fmt.Println("net dial err:", err)
		return nil
	}
	client.conn = conn

	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set your server ip (127.0.0.1 for default)")
	flag.IntVar(&serverPort, "port", 8888, "set your port(8888 for default)")

}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.group chating")
	fmt.Println("2.pair chating")
	fmt.Println("3.change username")
	fmt.Println("0.exit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please input valid integer")
		return false
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client conn write err:", err)
		return
	}

}

func (client *Client) PairChat() {
	var remoteName string
	var chatmsg string
	client.SelectUsers()
	fmt.Println("please input other name:")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println("please input your content:")
		fmt.Scanln(&chatmsg)
		for chatmsg != "exit" {
			if len(chatmsg) != 0 {
				sendmes := "to|" + remoteName + "|" + chatmsg + "\n"
				_, err := client.conn.Write([]byte(sendmes))
				if err != nil {
					fmt.Println("cilent conn write err:", err)
					break
				}
				chatmsg = ""
				fmt.Println("please input your message")
				fmt.Scanln(&chatmsg)
			}
		}
		client.SelectUsers()
		fmt.Println("please input other name:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		switch client.flag {
		case 1:
			client.GroupChat()
			break
		case 2:
			client.PairChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func (client *Client) GroupChat() {
	var sendmes string
	fmt.Println("please input your message")
	fmt.Scanln(&sendmes)
	for sendmes != "exit" {
		if len(sendmes) != 0 {
			sendmes = sendmes + "\n"
			_, err := client.conn.Write([]byte(sendmes))
			if err != nil {
				fmt.Println("cilent conn write err:", err)
				break
			}
			sendmes = ""
			fmt.Println("please input your message")
			fmt.Scanln(&sendmes)
		}
	}
}

func (this *Client) UpdateName() bool {
	fmt.Println("please input your name:")
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name + "\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("connection fail ")
	} else {
		fmt.Println("connection success")
	}

	go client.DealResponse()
	client.Run()
}
