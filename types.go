package ssm

import (
	"context"
	"sync"
	"time"
)

type (
	Machine struct {
		cbp *contextable_bp
	}

	Caller struct {
		continue_sig chan struct{}
	}

	State func(c *Caller) error

	contextable struct {
		ctx context.Context
	}

	boiler_plate struct {
		states      []State
		ctx         *contextable
		loop_tm     *time.Duration
		err_handler func(err error)
		core_lake   map[int]*data_lake[generic] // - unused now, for shared data
		last_state  int
	}

	generic interface{}

	data_lake[T generic] struct {
		lake sync.Pool
	}

	contextable_bp struct {
		*boiler_plate
	}
)
