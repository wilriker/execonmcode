package execonmcode

type Settings struct {
	SocketPath       string
	InterceptionMode string
	MCodes           MCodes
	Commands         Commands
	NoFlush          bool
	ExecAsync        bool
	Debug            bool
	Trace            bool
}
