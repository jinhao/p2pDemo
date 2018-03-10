package main

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jinhao/p2pDemo/proto"
)

var clientList *list.List
var clientMap map[string]*list.Element
var serverAddr = flag.String("addr", ":8080", "ip:port")

func main() {
	flag.Parse()
	serverAddr, err := net.ResolveUDPAddr("udp4", *serverAddr)
	if err != nil {
		log.Print(err)
		return
	}

	sk, err := net.ListenUDP("udp4", serverAddr)
	if err != nil {
		log.Print(err)
		return
	}
	defer sk.Close()

	// 初始化
	clientList = list.New()
	clientMap = make(map[string]*list.Element, 100)

	for {
		data := make([]byte, 4096)
		read, addr, err := sk.ReadFromUDP(data)
		if err != nil {
			log.Println("read failed!", err)
			continue
		}
		log.Printf("UDP: %d, %s\n", read, addr)

		// 解析协议
		var p proto.Proto
		err = json.Unmarshal(data[:read], &p)
		if err != nil {
			log.Print("json Unmarshal:", err)
			continue
		}
		switch p.Cmd {
		case proto.CMD_LOGIN:
			Login(addr.String())
		case proto.CMD_LOGOUT:
			Logout(addr.String())
		case proto.CMD_LIST:
			p := GetClients()
			p.Cmd = proto.CMD_LIST
			pData, err := json.Marshal(p)
			if err != nil {
				log.Print("json Marshal:", err)
			}
			sk.WriteToUDP(pData, addr)
		case proto.CMD_CONE:
			doCone(sk, p, addr.String())
		}
	}
}

// Login register to server
func Login(addr string) {
	log.Print("Login")
	// 判断client列表中是否有该客户端
	_, ok := clientMap[addr]
	if !ok {
		log.Printf("new UDP client: %s\n", addr)
		e := clientList.PushBack(addr)
		clientMap[addr] = e
	}
}

// Logout  .
func Logout(addr string) {
	log.Print("Logout")
	e, ok := clientMap[addr]
	if ok {
		v := clientList.Remove(e)
		log.Printf("Logout | %v", v)
		delete(clientMap, addr)
	}
}

func doCone(conn *net.UDPConn, p proto.Proto, fromAddr string) {
	fmt.Printf("fromAddr:%s coneAddr:%s", fromAddr, p.ConeAddr)
	coneAddr := p.ConeAddr
	p.Cmd = proto.CMD_CONE
	p.ConeAddr = fromAddr
	pData, err := json.Marshal(p)
	if err != nil {
		log.Print("json Marshal:", err)
	}
	coneAddrUDP, _ := net.ResolveUDPAddr("udp4", coneAddr)
	conn.WriteToUDP(pData, coneAddrUDP)

}

// GetClients get all register clients and save in a list
func GetClients() proto.Proto {
	log.Print("GetClients")
	p := proto.Proto{}
	for e := clientList.Front(); e != nil; e = e.Next() {
		val := e.Value.(string)
		p.Clients = append(p.Clients, val)
	}
	return p
}
