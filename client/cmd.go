package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"p2p/proto"
	"strconv"
)

var cmdlist = `
login: Login to server\n
list: list all clients\n
cone: cone to a clients\n
stop: stop and exit\n
help: show help info\n
server: reset new server addr\n
`

func ReadCmd(cmd chan CmdChan) {
	running := true
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(cmdlist)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch command {
		case "login":
			fmt.Println("cmd login")
			cmd <- CmdChan{cmd: proto.CMD_LOGIN}
		case "list":
			fmt.Println("cmd list")
			cmd <- CmdChan{cmd: proto.CMD_LIST}
		case "cone":
			fmt.Println("cmd cone")
			fmt.Printf("select clients to cone:")
			data, _, _ := reader.ReadLine()
			data_index, _ := strconv.Atoi(string(data))
			fmt.Printf("index:%d client num:%d\n", data_index, len(clientMap))
			if data_index < len(clientMap) {
				fmt.Printf("cone to:%s\n", clientMap[data_index])
				cmd <- CmdChan{cmd: proto.CMD_CONE, addr: clientMap[data_index]}
			}
		case "server":
			fmt.Println("cmd server")
			fmt.Printf("new server:")
			data, _, _ := reader.ReadLine()
			svr, err := net.ResolveUDPAddr("udp4", string(data))
			if err != nil {
				log.Printf("ResolverUdpAddr:%s err:%s", string(data), err.Error())
			} else {
				serverAddr = svr
			}
		case "stop":
			fmt.Println("cmd stop")
			os.Exit(0)
		case "help":
			fmt.Printf(cmdlist)
		}
	}
}
