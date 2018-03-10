# p2pDemo
简单测试UDP打洞用（待完善）。


# Build

	go get github.com/jinhao/p2pDemo
	make

# Usage

Client:

	p2pDemo git:(master) ✗ ./build/client -h
	Usage of ./build/client:
	  -addr string
		  bind addr (default ":8001")
	  -sAddr string
		  server addr (default "127.0.0.1:8080")

Server:

	p2pDemo git:(master) ✗ ./build/server -h
	Usage of ./build/server:
	  -addr string
			ip:port (default ":8080")

# RUN

1. 服务端放在具有公网ip1地址的服务器上，直接启动
	
	./server

2. 分别在两台内网机器上启动两个客户端C1,C2

	./client -sAddr=ip1:8080

3. C1、C2分别执行login命令，然后list查看已经login的客户端外网ip,如cIP1， cIP2

		cmd list
		2018/03/09 15:20:06 read data, remote addr:117.121.50.215:8080 3
		cmd_list len:2.
		0: 60.168.222.164:8001
		1: 112.64.215.22:61950

4. 从C1或C2端执行cone命令，根据提示，输入对端IP号，如输入1，此时已经两边开始打洞；

5. 观察打洞结果，如两边都打洞成功，会在屏幕上不断刷出对端发来的消息。

		2018/03/09 15:21:27 read data, remote addr:60.168.222.164:8001 5
		2018/03/09 15:21:27 recv msg:hello, I am 0.0.0.0:8001
		2018/03/09 15:21:27 read data, remote addr:117.121.50.215:8080 4
		2018/03/09 15:21:27 cone: 60.168.222.164:8001
		2018/03/09 15:21:32 read data, remote addr:60.168.222.164:8001 5
		2018/03/09 15:21:32 recv msg:hello, I am 0.0.0.0:8001
		2018/03/09 15:21:37 read data, remote addr:60.168.222.164:8001 5
