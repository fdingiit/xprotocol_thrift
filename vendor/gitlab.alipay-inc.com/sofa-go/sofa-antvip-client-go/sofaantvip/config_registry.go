package sofaantvip

import (
	"time"
)

const (
	DefaultRegistryLocatorTimeout = 5 * time.Second
)

type RegistryLocatorConfig struct {
	timeout time.Duration
}

func (rl *RegistryLocatorConfig) SetTimeout(t time.Duration) {
	rl.timeout = t
}
