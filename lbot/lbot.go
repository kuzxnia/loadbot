package lbot

type Lbot struct {
	config *Config
}

func NewLbot(config *Config) *Lbot {
	return &Lbot{
		config: config,
	}
}

func (l *Lbot) SetConfig(config *Config) {
	l.config = config
}
