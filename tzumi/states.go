package tzumi

import (
	"log"
)

type State uint8

const DISCONNECTED State = 0
const CONNECTED State = 1
const TUNING State = 2
const TUNED State = 3

func (s State) String() string {
	switch s {
	case DISCONNECTED:
		return "Disconnected"
	case CONNECTED:
		return "Connected"
	case TUNING:
		return "Tuning"
	case TUNED:
		return "Tuned"
	default:
		return "Unknown"
	}
}

func (t *TzumiMagicTV) setState(newState State) {
	if t.Debug {
		log.Printf("[DEBUG] New state is %s", newState.String())
	}
	t.State = newState
}
