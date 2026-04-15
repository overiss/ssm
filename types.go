package ssm

import (
	"sync"
	"time"
)

const (
	default_states_cap = 5
)

type (
	Machine struct {
		bps         []*boiler_plate
		replication *replication
	}

	Caller struct {
		continue_sig chan struct{}
		change_sig   chan string
		thread_id    int
	}

	State func(c *Caller) error

	Config struct {
		Loop_tm       time.Duration
		Err_handler   func(err error)
		Start_handler func(s *StartArg)
		Threads       int
	}

	StartArg struct {
		state_name string
		thread_id  int
	}

	replication struct {
		wg       *sync.WaitGroup
		replicas int
	}

	boiler_plate struct {
		states     []State
		names      []string
		state_map  map[string]int
		cfg        *Config
		last_state int
	}
)
