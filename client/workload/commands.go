package workload

// aggregate of commands

// start workload
// stop workload
// watch workload (in api will be run after start command when -w flag will be provided)

const (
	Start          = "start"
	Stop           = "stop"
	Watch          = "watch"
	GenerateConfig = "generate_config"
)

var Commands = []string{Start, Stop, Watch, GenerateConfig}

var nameToCommand = map[string]Command{
	Start:          NewStartCommand(),
	Stop:           &StartCommand{},
	Watch:          &WatchCommand{},
	GenerateConfig: &GenerateConfigCommand{},
}

type Command interface {
	handle() error
}

type StartCommand struct{}

func NewStartCommand() *StartCommand { return &StartCommand{} }

func (c *StartCommand) handle() error {
	// if watch arg - run watch
	return nil
}

type StopCommand struct{}

func (c *StopCommand) handle() error { return nil }

type WatchCommand struct{}

func (c *WatchCommand) handle() error { return nil }

type GenerateConfigCommand struct{}

func (c *GenerateConfigCommand) handle() error { return nil }
