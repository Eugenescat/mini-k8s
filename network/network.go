package network

import (
	"fmt"
	"net"
)

// StartTCPProxy 启动本地端口转发，将请求轮询转发到后端 podIP:podPort
func StartTCPProxy(localPort int, backends []string) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return err
	}
	fmt.Printf("[Network] 端口转发已启动，监听本地 %d，后端: %v\n", localPort, backends)
	go func() {
		var idx int
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			// 轮询选择后端
			target := backends[idx%len(backends)]
			idx++
			go handleProxy(conn, target)
		}
	}()
	return nil
}

func handleProxy(client net.Conn, backend string) {
	server, err := net.Dial("tcp", backend)
	if err != nil {
		client.Close()
		return
	}
	go func() { defer client.Close(); defer server.Close(); ioCopy(server, client) }()
	go func() { defer client.Close(); defer server.Close(); ioCopy(client, server) }()
}

func ioCopy(dst, src net.Conn) {
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, err2 := dst.Write(buf[:n]); err2 != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}
