package main

import (
	"flag"
	"log"
	"os"

	"github.com/Duet3D/DSF-APIs/godsfapi/v3/connection"
	"github.com/wilriker/execonmcode"
)

const (
	version = "5.1"
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
	version := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *version {
		log.Println(version)
		os.Exit(0)
	}

	if s.mCodes.Len() != s.commands.Len() {
		log.Fatal("Unequal amount of M-codes and commands given")
	}

	e := execonmcode.NewExecutor(s.socketPath, s.commands, s.mCodes, s.debug, s.trace)
	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
