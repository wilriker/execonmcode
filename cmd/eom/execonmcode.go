package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/Duet3D/DSF-APIs/godsfapi/v3/connection"
	"github.com/Duet3D/DSF-APIs/godsfapi/v3/connection/initmessages"
	"github.com/wilriker/execonmcode"
)

const (
	version = "5.2"
)

func main() {

	s := execonmcode.Settings{}

	flag.StringVar(&s.SocketPath, "socketPath", connection.FullSocketPath, "Path to socket")
	flag.StringVar(&s.InterceptionMode, "interceptionMode", string(initmessages.InterceptionModePre), "Interception mode to use")
	flag.Var(&s.MCodes, "mCode", "Code that will initiate execution of the command. This can be specified multiple times.")
	flag.Var(&s.Commands, "command", "Command to execute. This can be specified multiple times.")
	flag.BoolVar(&s.NoFlush, "noFlush", false, "Do not flush the code channel before executing the associated command")
	flag.BoolVar(&s.ExecAsync, "execAsync", false, "Run command to execute async and return success to DCS immediately")
	flag.BoolVar(&s.Debug, "debug", false, "Print debug output")
	flag.BoolVar(&s.Trace, "trace", false, "Print underlying requests/responses")
	version := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *version {
		log.Println(version)
		os.Exit(0)
	}

	if s.MCodes.Len() != s.Commands.Len() {
		log.Fatal("Unequal amount of M-codes and commands given")
	}

	switch strings.ToLower(s.InterceptionMode) {
	case strings.ToLower(string(initmessages.InterceptionModePre)):
		s.InterceptionMode = string(initmessages.InterceptionModePre)
	case strings.ToLower(string(initmessages.InterceptionModePost)):
		s.InterceptionMode = string(initmessages.InterceptionModePost)
	case strings.ToLower(string(initmessages.InterceptionModeExecuted)):
		s.InterceptionMode = string(initmessages.InterceptionModeExecuted)
	default:
		log.Fatal("Unsupported InterceptionMode", s.InterceptionMode)
	}

	e := execonmcode.NewExecutor(s)
	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
