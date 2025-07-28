package main

import (
	"net"
	"strconv"

	"github.com/alfredoprograma/tchat/internal/log"
)

type Args struct {
	Host net.IP
	Port int
}

func readArgs(args []string) Args {
	if len(args) < 3 {
		log.Log(log.LOG_LEVEL_FATAL, "expected host and port arguments")
	}

	rawHost := args[1]
	if host := net.ParseIP(rawHost); host == nil {
		log.Log(log.LOG_LEVEL_FATAL, "host argument should be a valid IPv4 address")
	}

	rawPort := args[2]
	port, err := strconv.Atoi(rawPort)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, "port argument should be an integer")
	}

	return Args{
		Host: net.ParseIP(args[1]),
		Port: port,
	}
}
