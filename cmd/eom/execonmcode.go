package main

import (
	"flag"
	"log"

	"github.com/wilriker/execonmcode"
	"github.com/wilriker/goduetapiclient/connection"
)

type settings struct {
	socketPath string
	mCode      int64
	command    string
}

func main() {
	s := settings{}

	flag.StringVar(&s.socketPath, "socketPath", connection.DefaultSocketPath, "Path to socket")
	flag.Int64Var(&s.mCode, "mCode", 7722, "Code that will initiate execution of the command")
	flag.StringVar(&s.command, "command", "", "Command to execute")
	flag.Parse()

	if s.mCode < 0 {
		log.Fatal("--mCode must be >= 0")
	}

	if s.command == "" {
		log.Fatal("--command must not be empty")
	}

	e := execonmcode.NewExecutor(s.socketPath, s.command, s.mCode)
	e.Run()
}
