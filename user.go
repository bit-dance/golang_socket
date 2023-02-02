package main

import (
	_"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn,server *Server) *User{
	userAddr:=conn.RemoteAddr().String()
	user:=&User{
		Name:userAddr,
		Addr:userAddr,
		C :make(chan string),
		conn:conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func(this *User) Online(){
	this.server.maplock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.maplock.Unlock()

	this.server.Broadcast(this,"Admining .....")
}

func(this *User) Offline(){
	this.server.maplock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.maplock.Unlock()

	this.server.Broadcast(this,"leaving .....")
}

func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

func (this *User) PassMessage(msg string){
	if msg=="who"{
		this.server.maplock.Lock()
		for _,cli:=range this.server.OnlineMap{
			onlinemsg:="["+cli.Addr+"]"+cli.Name+": online\n"
			this.SendMsg(onlinemsg)
		}
		this.server.maplock.Unlock()
	}else if len(msg)>7&&msg[:7]=="rename|"{
		newName:=strings.Split(msg,"|")[1]
		_,ok:=this.server.OnlineMap[newName]
		if ok{
			this.SendMsg("This name has been used!\n")
		}else{
			this.server.maplock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName]=this
			this.server.maplock.Unlock()

			this.Name=newName
			this.SendMsg("You have update name :"+newName+"\n")
		}
		
	}else if len(msg)>4&&msg[:3]=="to|"{
		remoteName:=strings.Split(msg,"|")[1]
		if remoteName==""{
			this.SendMsg("Please correct your format like this \"to|jack|hello\"\n")
			return
		}
		remoteUser,ok:=this.server.OnlineMap[remoteName]

		if !ok{
			this.SendMsg("Please input correct user name!\n")
			return 
		}

		content:=strings.Split(msg,"|")[2]
		if content==""{
			this.SendMsg("Invalid message,please input some words!!\n")
			return 
		}
		remoteUser.SendMsg(this.Name+":   "+content+"\n")
		

	}else {
		this.server.Broadcast(this,msg)
	}
}


func(this *User) ListenMessage(){
	for{
		msg:=<-this.C
		this.conn.Write([]byte(msg+"\n"))
	}
}