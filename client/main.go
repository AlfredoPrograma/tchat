package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/alfredoprograma/tchat/internal/events"
	"github.com/alfredoprograma/tchat/internal/log"
)

func NewConn(connSettings connSettings) *net.TCPConn {
	conn, err := net.DialTCP("tcp", nil, connSettings.addr)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, fmt.Sprintf("cannot connect to host %s:%d", connSettings.addr.IP, connSettings.addr.Port))
	}

	event := events.NewRegisterUserEvent(connSettings.username)
	serialized := events.Serialize(event)
	_, err = conn.Write(serialized)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, fmt.Sprintf("cannot register with username %s", connSettings.username))
	}

	return conn
}

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	rw := bufio.NewReadWriter(r, w)
	settings := handleSettingsPrompt(rw)
	conn := NewConn(settings)
	end := make(chan bool)

	go handleMessagePrompt(rw, conn)
	<-end
}
