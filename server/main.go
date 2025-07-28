package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

const BUFFER_SIZE = 512

func NewAddr(port int) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	}
}

func NewListener(addr *net.TCPAddr) *net.TCPListener {
	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log(LOG_LEVEL_FATAL, fmt.Sprintf("cannot start listener at addr %s", addr))
	}

	return listener
}

func receiveConnections(listener *net.TCPListener, eventsCh chan Event) {
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log(LOG_LEVEL_ERROR, fmt.Sprintf("cannot handle connection from %s", conn.RemoteAddr().String()))
			continue
		}

		go handleConnection(conn, eventsCh)
	}
}

func handleConnection(conn *net.TCPConn, eventsChan chan Event) {
	for {
		buf := make([]byte, 0, BUFFER_SIZE)
		_, err := conn.Read(buf)

		if err != nil {
			log(LOG_LEVEL_ERROR, fmt.Sprintf("cannot read from connection %s", conn.RemoteAddr().String()))
			continue
		}

		var incomingEvent Event

		if err := json.Unmarshal(buf, &incomingEvent); err != nil {
			log(LOG_LEVEL_ERROR, fmt.Sprintf("cannot parse json message from connection %s", conn.RemoteAddr().String()))
			continue
		}

		eventsChan <- incomingEvent
	}
}

func handleEvents(eventsCh chan Event) {
	for event := range eventsCh {
		log(LOG_LEVEL_INFO, fmt.Sprintf("handling %s event", event.Kind))
	}
}

type Chat struct {
	mu    sync.Mutex
	conns []net.TCPConn
}

func NewChat() Chat {
	return Chat{
		conns: make([]net.TCPConn, 0),
	}
}

func main() {
	args := readArgs(os.Args)
	addr := NewAddr(args.Port)
	listener := NewListener(addr)
	eventsCh := make(chan Event)
	end := make(chan bool)

	go receiveConnections(listener, eventsCh)
	go handleEvents(eventsCh)

	<-end
}
