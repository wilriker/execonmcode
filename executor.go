package execonmcode

import (
	"log"
	"os/exec"
	"strings"

	"github.com/wilriker/goduetapiclient/connection"
	"github.com/wilriker/goduetapiclient/connection/initmessages"
	"github.com/wilriker/goduetapiclient/types"
)

type Executor struct {
	socketPath string
	mCode      int64
	command    string
	args       []string
}

func NewExecutor(socketPath, command string, mCode int64) *Executor {
	s := strings.Split(command, " ")
	a := []string{}
	if len(s) > 1 {
		a = s[1:]
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
			cmd := exec.Command(e.command, e.args...)
			err := cmd.Run()
			if err != nil {
				err = ic.ResolveCode(types.Error, err.Error())
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
