package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/alfredoprograma/tchat/internal/events"
	"github.com/alfredoprograma/tchat/internal/log"
)

func NewAddr(host net.IP, port int) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   host,
		Port: port,
	}
}

func NewConn(addr *net.TCPAddr) *net.TCPConn {
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, fmt.Sprintf("cannot connect to host %s:%d", addr.IP, addr.Port))
	}

	return conn
}

func main() {
	args := readArgs(os.Args)
	addr := NewAddr(args.Host, args.Port)
	conn := NewConn(addr)

	myEvent := events.Event{
		Kind: "REGISTER",
		Payload: map[string]string{
			"username": "alf2001",
		},
	}

	buf, err := json.Marshal(myEvent)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, "cannot parse event")
	}

	_, err = conn.Write(buf)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, "cannot write to conn")
	}
}
