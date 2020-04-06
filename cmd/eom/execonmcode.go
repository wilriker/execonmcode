package main

import (
	"flag"
	"log"

	"github.com/Duet3D/DSF-APIs/godsfapi/connection"
	"github.com/wilriker/execonmcode"
)

type settings struct {
	socketPath string
	mCodes     execonmcode.MCodes
	commands   execonmcode.Commands
	debug      bool
	trace      bool
}

func main() {

	s := settings{}

	flag.StringVar(&s.socketPath, "socketPath", connection.FullSocketPath, "Path to socket")
	flag.Var(&s.mCodes, "mCode", "Code that will initiate execution of the command")
	flag.Var(&s.commands, "command", "Command to execute")
	flag.BoolVar(&s.debug, "debug", false, "Print debug output")
	flag.BoolVar(&s.trace, "trace", false, "Print underlying requests/responses")
	flag.Parse()

	if s.mCodes.Len() != s.commands.Len() {
		log.Fatal("Unequal amount of M-codes and commands given")
	}

	e := execonmcode.NewExecutor(s.socketPath, s.commands, s.mCodes, s.debug, s.trace)
	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
