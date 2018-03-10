package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jinhao/p2pDemo/proto"
)

func doLogin(conn *net.UDPConn, serverAddr *net.UDPAddr) error {
	// Login
	p := proto.Proto{}
	p.Cmd = proto.CMD_LOGIN
	sData, err := json.Marshal(p)
	if err != nil {
		log.Print(err)
		return err
	}
	return WriteData(conn, sData, serverAddr)
}

func listClient(conn *net.UDPConn, serverAddr *net.UDPAddr) error {
	p := proto.Proto{}
	p.Cmd = proto.CMD_LIST
	sData, err := json.Marshal(p)
	if err != nil {
		log.Print(err)
		return err
	}

	return WriteData(conn, sData, serverAddr)
}

func sendConeReq(conn *net.UDPConn, serverAddr *net.UDPAddr, coneAddr string) error {
	p := proto.Proto{}
	p.Cmd = proto.CMD_CONE
	p.ConeAddr = coneAddr
	sData, err := json.Marshal(p)
	if err != nil {
		log.Print(err)
		return err
	}

	WriteData(conn, sData, serverAddr)
	go doCone(conn, coneAddr)

	return nil
}

// WriteData .
func WriteData(conn *net.UDPConn, data []byte, serverAddr *net.UDPAddr) error {
	_, err := conn.WriteToUDP(data, serverAddr)
	if err != nil {
		log.Print(err)
	}
	return err
}

// ReadData .
func ReadData(conn *net.UDPConn) error {
	for {
		data := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Print(err)
			continue
		}
		p := proto.Proto{}
		err = json.Unmarshal(data[:n], &p)
		if err != nil {
			log.Printf("json Unmarshal err:%s", err)
		}

		log.Printf("read data, remote addr:%s %d", addr, p.Cmd)

		switch p.Cmd {
		case proto.CMD_LIST:
			fmt.Printf("cmd_list len:%d.\n", len(p.Clients))
			clientMap = make(map[int]string)
			for i, v := range p.Clients {
				fmt.Printf("%d: %s\n", i, v)
				clientMap[i] = v
			}
		case proto.CMD_CONE:
			log.Printf("cone: %s\n", p.ConeAddr)
			go doCone(conn, p.ConeAddr)

		case proto.CMD_MSG:
			log.Printf("recv msg:%s\n", p.Data)
		}
	}

}

func doCone(localUDPConn *net.UDPConn, addr string) error {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Printf("doCone | ResolveUDPAddr err:%s\n", err.Error())
		return err
	}
	for {
		p := proto.Proto{}
		p.Cmd = proto.CMD_MSG
		p.Data = fmt.Sprintf("hello, I am %s", localUDPConn.LocalAddr().String())
		msg, _ := json.Marshal(&p)
		localUDPConn.WriteToUDP(msg, remoteUDPAddr)
		time.Sleep(5 * time.Second)
	}

}
