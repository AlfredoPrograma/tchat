package main

import (
	"net"
	"sync"
)

type User struct {
	username string
	conn     *net.TCPConn
}

type Chat struct {
	mu    sync.Mutex
	conns map[string]User
}

func NewChat() *Chat {
	return &Chat{
		conns: make(map[string]User, 0),
	}
}


