package execonmcode

import (
	"fmt"
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
	mCode      int64
	command    string
	args       []string
}

func NewExecutor(socketPath, command string, mCode int64) *Executor {
	s := strings.Split(command, " ")
	a := make([]string, 0)
	if len(s) > 1 {
		for _, p := range s[1:] {
			pp := strings.TrimSpace(p)
			if pp != "" {
				a = append(a, p)
			}
		}
	}
	c := s[0]
	return &Executor{
		socketPath: socketPath,
		command:    c,
		args:       a,
		mCode:      mCode,
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
			log.Println("Error:", err)
			continue
		}
		if c.Type == types.MCode && c.MajorNumber != nil && *c.MajorNumber == e.mCode {
			cmd := exec.Command(e.command, e.getArgs(c)...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				err = ic.ResolveCode(types.Error, fmt.Sprintf("%s: %s", err.Error(), string(output)))
			} else {
				err = ic.ResolveCode(types.Success, "")
			}
			if err != nil {
				log.Println("Error:", err)
			}
		} else {
			ic.IgnoreCode()
		}
	}
}

func (e *Executor) getArgs(c *commands.Code) []string {
	args := make([]string, 0)
	for _, v := range e.args {
		if strings.HasPrefix(v, variablePrefix) {
			vl := strings.TrimSpace(strings.ToUpper(strings.TrimLeft(v, variablePrefix)))
			if len(vl) == 1 {
				if pv := c.Parameter(vl); pv != nil {
					v = pv.AsString()
				}
			}
		}
		args = append(args, v)
	}
	return args
}
