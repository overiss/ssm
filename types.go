package ssm

import (
	"context"
	"time"
)

type (
	Machine struct {bp *boiler_plate}

	Caller struct {continue_sig chan struct{}}

	State func(c *Caller) error

	Config struct {
		Loop_tm       time.Duration
		Err_handler   func(err error)
		Start_handler func(state_name string)
	}

	boiler_plate struct {
		states     []State
		names      []string
		ctx        context.Context
		cfg        *Config
		last_state int
	}
)
