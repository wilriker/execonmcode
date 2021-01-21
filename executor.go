package execonmcode

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/Duet3D/DSF-APIs/godsfapi/v3/commands"
	"github.com/Duet3D/DSF-APIs/godsfapi/v3/connection"
	"github.com/Duet3D/DSF-APIs/godsfapi/v3/connection/initmessages"
	"github.com/Duet3D/DSF-APIs/godsfapi/v3/machine/messages"
)

const (
	variablePrefix = "%"
)

type Executor struct {
	socketPath   string
	mode         initmessages.InterceptionMode
	mCodes       map[int64]int
	commands     Commands
	execAsync    bool
	returnOutput bool
	flush        bool
	debug        bool
	trace        bool
}

func NewExecutor(s Settings) *Executor {
	mc := make(map[int64]int)
	for i, m := range s.MCodes {
		mc[m] = i
		if s.Debug {
			cmd, args, err := s.Commands.Get(i)
			if err != nil {
				log.Println(m, err)
			}
			log.Printf("%d: %s %s", m, cmd, strings.Join(args, " "))
		}
	}
	return &Executor{
		socketPath:   s.SocketPath,
		mode:         initmessages.InterceptionMode(s.InterceptionMode),
		mCodes:       mc,
		commands:     s.Commands,
		execAsync:    s.ExecAsync,
		returnOutput: s.ReturnOutput,
		flush:        !s.NoFlush,
		debug:        s.Debug,
		trace:        s.Trace,
	}
}

func (e *Executor) Run() error {

	ic := connection.InterceptConnection{}
	ic.Debug = e.trace
	err := ic.Connect(e.mode, e.socketPath)
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
		if c.Type == commands.MCode && c.MajorNumber != nil {
			i, ok := e.mCodes[*c.MajorNumber]
			if !ok {
				ic.IgnoreCode()
				continue
			}
			if e.flush {
				success, err := ic.Flush(c.Channel)
				if !success || err != nil {
					log.Println("Could not Flush. Cancelling code")
					ic.CancelCode()
					continue
				}
			}
			comd, a, err := e.commands.Get(i)
			if err != nil {
				ic.ResolveCode(messages.Error, err.Error())
			} else {
				cmd := exec.Command(comd, e.getArgs(c, a)...)
				if e.debug {
					log.Println("Executing:", cmd)
				}

				// If we should exec async run it as goroutine and return success
				if e.execAsync {
					go cmd.Run()
					err = ic.ResolveCode(messages.Success, "")
				} else {
					output, err := cmd.CombinedOutput()
					if err != nil {
						err = ic.ResolveCode(messages.Error, fmt.Sprintf("%s: %s", err.Error(), string(output)))
					} else {
						msg := ""
						if e.returnOutput {
							msg = string(output)
						}
						err = ic.ResolveCode(messages.Success, msg)
					}
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
