package execonmcode

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/wilriker/goduetapiclient/commands"
	"github.com/wilriker/goduetapiclient/connection"
	"github.com/wilriker/goduetapiclient/connection/initmessages"
	"github.com/wilriker/goduetapiclient/types"
)

const (
	variablePrefix = "%"
)

type Executor struct {
	socketPath string
	mCodes     map[int64]int
	commands   Commands
	debug      bool
}

func NewExecutor(socketPath string, commands Commands, mCodes MCodes, debug bool) *Executor {
	mc := make(map[int64]int)
	for i, m := range mCodes {
		mc[m] = i
		if debug {
			cmd, args, err := commands.Get(i)
			if err != nil {
				log.Println(m, err)
			}
			log.Printf("%d: %s %s", m, cmd, strings.Join(args, " "))
		}
	}
	return &Executor{
		socketPath: socketPath,
		mCodes:     mc,
		commands:   commands,
		debug:      debug,
	}
}

func (e *Executor) Run() {

	ic := connection.InterceptConnection{}
	err := ic.Connect(initmessages.InterceptionModePre, e.socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer ic.Close()

	for {
		c, err := ic.ReceiveCode()
		if err != nil {
			if err == io.EOF {
				log.Println("Connection to DCS closed")
				break
			}
			log.Printf("Error receiving code: %s", err)
			continue
		}
		if c.Type == types.MCode && c.MajorNumber != nil {
			i, ok := e.mCodes[*c.MajorNumber]
			if !ok {
				ic.IgnoreCode()
				continue
			}
			comd, a, err := e.commands.Get(i)
			if err != nil {
				ic.ResolveCode(types.Error, err.Error())
			} else {
				cmd := exec.Command(comd, e.getArgs(c, a)...)
				if e.debug {
					log.Println("Executing:", cmd)
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					err = ic.ResolveCode(types.Error, fmt.Sprintf("%s: %s", err.Error(), string(output)))
				} else {
					err = ic.ResolveCode(types.Success, "")
				}
				if err != nil {
					log.Println("Error executing command:", err)
				}
			}
		} else {
			ic.IgnoreCode()
		}
	}
}

func (e *Executor) getArgs(c *commands.Code, args []string) []string {
	a := make([]string, 0)
	for _, v := range args {
		if strings.HasPrefix(v, variablePrefix) {
			vl := strings.TrimSpace(strings.ToUpper(strings.TrimLeft(v, variablePrefix)))
			if len(vl) == 1 {
				if pv := c.Parameter(vl); pv != nil {
					v = pv.AsString()
				}
			}
		}
		a = append(a, v)
	}
	return a
}
