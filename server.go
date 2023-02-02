package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Port int
	OnlineMap map[string]*User
	maplock sync.RWMutex
	Message chan string
}

func NewServer(ip string,port int) *Server {
	server:=&Server {
		Ip:ip,
		Port:port,
		OnlineMap:make(map[string]*User),
		Message:make(chan string),
	}
	return server
}

func(this *Server) ListenMessageer(){
	for  {
		msg:=<-this.Message
		this.maplock.Lock()
		for _,cli:=range this.OnlineMap {
			cli.C <-msg
		}
		this.maplock.Unlock()
	}
}

func (this *Server) Broadcast(user *User,msg string){
	sendMsg:="["+user.Addr+"]"+user.Name+":"+msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("server listening.......")
	user:=NewUser(conn,this)
	user.Online()

	//make a channel to moniter the status
	Islive:=make(chan bool)
	go func(){
		buf:=make([]byte,4096)
		for{
			n,err:=conn.Read(buf)
			if n==0{
				user.Offline()
				return
			}
			if err!=nil && err!=io.EOF{
				fmt.Println("Conn Read err:",err)
				return 
			}
			msg:=string(buf[:n-1])
			user.PassMessage(msg)
			Islive<-true
		}
	}()
	
	for{
		select{
		case <-Islive:
			//do nothing to refresh the timer
		case <-time.After(time.Second*300):
			user.SendMsg("You are evicted!!!")
			delete(this.OnlineMap,user.Name)
			user.conn.Close()
			close(user.C)
			return 
		}
	}
	
}

func (this *Server) Start() {
	listener,err := net.Listen("tcp", fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err!=nil {
		fmt.Println("net listener err:",err)
		return
	}

	defer listener.Close()

	go this.ListenMessageer()

	for {
		conn,err := listener.Accept()
		if err!=nil{
			fmt.Println("listener accept err:",err)
			continue
		}

		go this.Handler(conn)
	}

}
