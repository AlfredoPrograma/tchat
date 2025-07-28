package main

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"

	"github.com/alfredoprograma/tchat/internal/events"
	"github.com/alfredoprograma/tchat/internal/log"
)

func handleHostPrompt(rw *bufio.ReadWriter) net.IP {
	rw.WriteString("> Host: ")
	rw.Flush()

	for {
		input, _ := rw.ReadString('\n')
		ip := net.ParseIP(strings.Trim(input, "\n"))
		if ip == nil {
			rw.WriteString("Invalid host; try another...\n")
			rw.Flush()
			continue
		}

		return ip
	}
}

func handlePortPrompt(rw *bufio.ReadWriter) int {
	rw.WriteString("> Port: ")
	rw.Flush()

	for {
		input, _ := rw.ReadString('\n')
		port, err := strconv.Atoi(strings.Trim(input, "\n"))

		if err != nil {
			rw.WriteString("Invalid port; try another...\n")
			rw.Flush()
			continue
		}

		return port
	}
}

func handleUsernamePrompt(rw *bufio.ReadWriter) string {
	rw.WriteString("> Username: ")
	rw.Flush()

	for {
		input, _ := rw.ReadString('\n')

		if input == "" {
			rw.WriteString("Invalid username; try another...\n")
			rw.Flush()
			continue
		}

		return strings.Trim(input, "\n")
	}
}

type connSettings struct {
	username string
	addr     *net.TCPAddr
}

func handleSettingsPrompt(rw *bufio.ReadWriter) connSettings {
	host := handleHostPrompt(rw)
	port := handlePortPrompt(rw)
	username := handleUsernamePrompt(rw)

	return connSettings{
		username: username,
		addr: &net.TCPAddr{
			IP:   host,
			Port: port,
		},
	}
}

func handleMessagePrompt(rw *bufio.ReadWriter, conn *net.TCPConn) {
	for {
		input, _ := rw.ReadString('\n')

		event := events.NewSendMessageEvent(input)
		raw, err := json.Marshal(event)

		if err != nil {
			log.Log(log.LOG_LEVEL_ERROR, "cannot send register user event")
			return
		}

		conn.Write(raw)
	}
}
