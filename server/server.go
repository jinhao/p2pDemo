package main

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"p2p/proto"
)

var clientList *list.List
var clientMap map[string]*list.Element
var serverAddr = flag.String("addr", ":8080")

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
			p_data, err := json.Marshal(p)
			if err != nil {
				log.Print("json Marshal:", err)
			}
			sk.WriteToUDP(p_data, addr)
		case proto.CMD_CONE:
			doCone(sk, p, addr.String())
		}
	}
}

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

func Logout(addr string) {
	log.Print("Logout")
	e, ok := clientMap[addr]
	if ok {
		v := clientList.Remove(e)
		log.Print("Logout | %v", v)
		delete(clientMap, addr)
	}
}

func doCone(conn *net.UDPConn, p proto.Proto, fromAddr string) {
	fmt.Printf("fromAddr:%s coneAddr:%s", fromAddr, p.ConeAddr)
	coneAddr := p.ConeAddr
	p.Cmd = proto.CMD_CONE
	p.ConeAddr = fromAddr
	p_data, err := json.Marshal(p)
	if err != nil {
		log.Print("json Marshal:", err)
	}
	coneAddrUdp, _ := net.ResolveUDPAddr("udp4", coneAddr)
	conn.WriteToUDP(p_data, coneAddrUdp)

}

func GetClients() proto.Proto {
	log.Print("GetClients")
	p := proto.Proto{}
	for e := clientList.Front(); e != nil; e = e.Next() {
		val := e.Value.(string)
		p.Clients = append(p.Clients, val)
	}
	return p
}
