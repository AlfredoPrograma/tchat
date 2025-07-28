package main

import (
	"strconv"

	"github.com/alfredoprograma/tchat/internal/log"
)

type Args struct {
	Port int
}

func readArgs(args []string) Args {
	if len(args) < 2 {
		log.Log(log.LOG_LEVEL_FATAL, "expected port argument")
	}

	rawPort := args[1]
	port, err := strconv.Atoi(rawPort)

	if err != nil {
		log.Log(log.LOG_LEVEL_FATAL, "port argument should be an integer")
	}

	return Args{Port: port}
}
