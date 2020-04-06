package execonmcode

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/Duet3D/DSF-APIs/godsfapi/commands"
	"github.com/Duet3D/DSF-APIs/godsfapi/connection"
	"github.com/Duet3D/DSF-APIs/godsfapi/connection/initmessages"
	"github.com/Duet3D/DSF-APIs/godsfapi/types"
)

const (
	variablePrefix = "%"
)

type Executor struct {
	socketPath string
	mCodes     map[int64]int
	commands   Commands
	debug      bool
	trace      bool
}

func NewExecutor(socketPath string, commands Commands, mCodes MCodes, debug, trace bool) *Executor {
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
		trace:      trace,
	}
}

func (e *Executor) Run() error {

	ic := connection.InterceptConnection{}
	ic.Debug = e.trace
	err := ic.Connect(initmessages.InterceptionModePre, e.socketPath)
	if err != nil {
		return err
	}
	defer ic.Close()

	for {
		c, err := ic.ReceiveCode()
		if err != nil {
			if err == io.EOF {
				log.Println("Connection to DCS closed")
				return err
			}
			if _, ok := err.(*connection.DecodeError); ok {
				// If it is "just" a problem with decoding ignore the received code
				// as it otherwise will block DCS
				ic.IgnoreCode()
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
