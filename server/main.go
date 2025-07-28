package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/alfredoprograma/tchat/internal/events"
	"github.com/alfredoprograma/tchat/internal/log"
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
		log.Log(log.LOG_LEVEL_FATAL, fmt.Sprintf("cannot start listener at addr %s", addr))
	}

	return listener
}

func receiveConnections(listener *net.TCPListener, eventsCh chan events.Event) {
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("cannot handle connection from %s", conn.RemoteAddr().String()))
			continue
		}

		go handleConnection(conn, eventsCh)
	}
}

func handleConnection(conn *net.TCPConn, eventsChan chan events.Event) {
	for {
		buf := make([]byte, BUFFER_SIZE)
		readLen, err := conn.Read(buf)

		if err != nil {
			// log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("cannot read from connection %s", conn.RemoteAddr().String()))
			continue
		}

		if readLen == 0 {
			continue
		}

		var incomingEvent events.Event

		if err := json.Unmarshal(buf[:readLen], &incomingEvent); err != nil {
			log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("cannot parse json message from connection %s", conn.RemoteAddr().String()))
			continue
		}

		incomingEvent.Meta = events.EventMetadata{
			Conn: conn,
		}

		eventsChan <- incomingEvent
	}
}

func registerUser(username string, conn *net.TCPConn, chat *Chat) {
	chat.mu.Lock()
	chat.conns[username] = conn
	chat.mu.Unlock()

	log.Log(log.LOG_LEVEL_INFO, fmt.Sprintf("user %s registered in chat room", username))
}

func handleEvents(eventsCh chan events.Event, chat *Chat) {
	for event := range eventsCh {
		log.Log(log.LOG_LEVEL_INFO, fmt.Sprintf("handling %s event", event.Kind))

		switch event.Kind {
		case events.REGISTER_USER_EVENT:
			// TODO: maybe add validations ???
			payload := event.Payload.(map[string]any)
			username := payload["Username"].(string)
			registerUser(username, event.Meta.Conn, chat)
		case events.SEND_MESSAGE_EVENT:
			// broadcastMessage()
		default:
			log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("invalid event kind %s", event.Kind))
		}
	}
}

type Chat struct {
	mu    sync.Mutex
	conns map[string]*net.TCPConn
}

func NewChat() *Chat {
	return &Chat{
		conns: make(map[string]*net.TCPConn, 0),
	}
}

func main() {
	args := readArgs(os.Args)
	addr := NewAddr(args.Port)
	listener := NewListener(addr)
	chat := NewChat()
	eventsCh := make(chan events.Event)
	end := make(chan bool)

	go receiveConnections(listener, eventsCh)
	go handleEvents(eventsCh, chat)

	<-end
}
