package main

import (
	"flag"
	"log"
	"os/exec"

	"github.com/wilriker/goduetapiclient/connection"
	"github.com/wilriker/goduetapiclient/connection/initmessages"
	"github.com/wilriker/goduetapiclient/types"
)

type settings struct {
	SocketPath    string
	ShutdownMCode int64
}

func main() {
	s := settings{}

	flag.StringVar(&s.SocketPath, "socketPath", connection.DefaultSocketPath, "Path to socket")
	flag.Int64Var(&s.ShutdownMCode, "shutdownMCode", 7722, "Code that will initiate shutdown of the SBC")
	flag.Parse()

	ic := connection.InterceptConnection{}
	err := ic.Connect(initmessages.InterceptionModePre, s.SocketPath)
	if err != nil {
		panic(err)
	}
	defer ic.Close()

	for {
		c, err := ic.ReceiveCode()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		if c.Type == types.MCode && c.MajorNumber != nil && *c.MajorNumber == s.ShutdownMCode {
			err = ic.ResolveCode(types.Success, "")
			if err != nil {
				log.Println("Error:", err)
			}
			cmd := exec.Command("poweroff")
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			ic.IgnoreCode()
		}
	}
}
