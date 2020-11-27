package stateconn

import "time"

type State uint8

const (
	StateNew State = iota
	StateActive
	StateIdle
	StateHijacked
	StateClosed
)

func packState(state State) int64 {
	return time.Now().Unix()<<8 | int64(state)
}

func (s State) String() string {
	switch s {
	case StateNew:
		return "new"
	case StateActive:
		return "active"
	case StateIdle:
		return "idle"
	case StateClosed:
		return "closed"
	case StateHijacked:
		return "hijacked"
	default:
		return "unknown state"
	}
}
