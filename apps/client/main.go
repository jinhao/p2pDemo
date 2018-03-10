package main

import (
	"flag"
	"log"
	"net"
	"p2p/proto"
)

var clientMap map[int]string

type CmdChan struct {
	cmd  int
	addr string
}

var lisAddr = flag.String("addr", ":8001", "bind addr")
var svrAddr = flag.String("sAddr", "127.0.0.1:8080", "sever addr")
var serverAddr *net.UDPAddr

func main() {
	flag.Parse()
	var err error
	//serverAddr, err := net.ResolveUDPAddr("udp4", "117.121.50.216:8080")
	serverAddr, err = net.ResolveUDPAddr("udp4", *svrAddr)
	if err != nil {
		log.Print(err)
		return
	}
	listenAddr, err := net.ResolveUDPAddr("udp4", *lisAddr)
	if err != nil {
		log.Print(err)
		return
	}

	sk, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		log.Print("connect fail", err)
		return
	}
	defer sk.Close()

	cmdChan := make(chan CmdChan)
	go ReadCmd(cmdChan)
	go ReadData(sk)

	for {
		command := <-cmdChan
		switch command.cmd {
		case proto.CMD_LOGIN:
			doLogin(sk, serverAddr)
		case proto.CMD_LIST:
			listClient(sk, serverAddr)
		case proto.CMD_CONE:
			log.Printf("cmd:cone, addr:%s", command.addr)
			sendConeReq(sk, serverAddr, command.addr)

		case proto.CMD_LOGOUT:

		default:
			log.Printf("unsupported command:%d", command)
		}
	}
}
