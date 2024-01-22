package workload

import (
	"fmt"

	"github.com/kuzxnia/loadcli/cli"
)

type Manager struct {
	command Command
	// context, args and more stuff
}

func NewManager(request cli.Request) (*Manager, error) {
	command, ok := nameToCommand[request.Command]

	if !ok {
		return nil, fmt.Errorf("command not found: %s", request.Command)
	}

	return &Manager{
		command: command,
	}, nil
}

func (m *Manager) Run() (err error) {
	err = m.command.handle()

	return
}
