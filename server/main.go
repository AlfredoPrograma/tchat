package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/alfredoprograma/tchat/internal/events"
	"github.com/alfredoprograma/tchat/internal/log"
)

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

func receiveConnections(eventsCh chan events.Event, chat *Chat, listener *net.TCPListener) {
	log.Log(log.LOG_LEVEL_INFO, fmt.Sprintf("started to accept TCP connections at %s", listener.Addr().String()))
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("cannot handle connection from %s", conn.RemoteAddr().String()))
			continue
		}

		go handleConnection(conn, chat, eventsCh)
	}
}

func handleConnection(conn *net.TCPConn, chat *Chat, eventsChan chan events.Event) {
	for {
		buf := make([]byte, events.BUFFER_SIZE)
		readLen, err := conn.Read(buf)

		if err != nil {
			if errors.Is(err, io.EOF) {
				removeUser(conn, chat)
				break
			}

			continue
		}

		if readLen == 0 {
			continue
		}

		incomingEvent := events.Deserialize(buf[:readLen])
		incomingEvent.Meta = events.EventMetadata{
			Conn: conn,
		}

		eventsChan <- incomingEvent
	}
}

func registerUser(username string, conn *net.TCPConn, chat *Chat) {
	chat.mu.Lock()
	chat.conns[conn.RemoteAddr().String()] = User{
		username: username,
		conn:     conn,
	}
	chat.mu.Unlock()

	log.Log(log.LOG_LEVEL_INFO, fmt.Sprintf("user %s registered in chat room", username))
}

func removeUser(conn *net.TCPConn, chat *Chat) {
	chat.mu.Lock()
	username := chat.conns[conn.RemoteAddr().String()].username
	delete(chat.conns, conn.RemoteAddr().String())
	log.Log(log.LOG_LEVEL_INFO, fmt.Sprintf("user %s left from chat room", username))
	chat.mu.Unlock()
}

func broadcast(message string, emitter *net.TCPConn, chat *Chat) {
	chat.mu.Lock()
	emitterUsername := chat.conns[emitter.RemoteAddr().String()].username
	buf := fmt.Sprintf("[%s]: %s", emitterUsername, message)
	log.Log(log.LOG_LEVEL_INFO, buf)

	for addr, user := range chat.conns {
		// Omit broadcast message to emitter
		if addr == emitter.RemoteAddr().String() {
			continue
		}

		user.conn.Write([]byte(buf))
	}

	chat.mu.Unlock()
}

func handleEvents(eventsCh chan events.Event, chat *Chat) {
	for event := range eventsCh {
		switch event.Kind {
		case events.REGISTER_USER_EVENT:
			payload := event.Payload.(events.RegisterUserPayload)
			registerUser(payload.Username, event.Meta.Conn, chat)
		case events.SEND_MESSAGE_EVENT:
			payload := event.Payload.(events.SendMessagePayload)
			broadcast(payload.Content, event.Meta.Conn, chat)
		default:
			log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("invalid event kind %s", event.Kind))
		}
	}
}

func main() {
	args := readArgs(os.Args)
	addr := NewAddr(args.Port)
	listener := NewListener(addr)
	eventsCh := make(chan events.Event)
	end := make(chan bool)
	chat := NewChat()

	go receiveConnections(eventsCh, chat, listener)
	go handleEvents(eventsCh, chat)

	<-end
}
