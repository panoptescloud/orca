package hostsys

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
}

type HostSystem struct {
	tui tui
}

func NewHostSystem(tui tui) *HostSystem {
	return &HostSystem{
		tui: tui,
	}
}
