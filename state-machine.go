package ssm

import (
	"context"
	"sync"
	"time"
)

// CreateMachine - creates a set of methods for configuring the state machine
func CreateMachine() *boiler_plate {
	return &boiler_plate{states: make([]State, 0, default_states_cap), cfg: new(Config), state_map: make(map[string]int)}
}

// ApplyCfg - applies configurations to the state machine
func (b *boiler_plate) ApplyCfg(cfg *Config) *boiler_plate { b.cfg = cfg; return b }

// AddState - adds an executable state to the machine. You can specify a custom
// state name as the second argument, otherwise it will be determined by default.
func (b *boiler_plate) AddState(s State, state_name string) *boiler_plate {
	b.states, b.names = append(b.states, s), append(b.names, state_name)
	b.state_map[state_name] = len(b.states) - 1
	return b
}

// Build - assembles the configured machine into an executable method
func (b *boiler_plate) Build() *Machine {
	rc := 1
	real_c := b.cfg.Threads
	if real_c > 1 {
		rc = real_c
	}
	new_wg := new(sync.WaitGroup)
	m := &Machine{bps: b.init_replicated_plates(rc), replication: &replication{wg: new_wg, replicas: rc}}
	return m
}

func (m *Machine) Run(ctx context.Context) {
	fn := func(gid int) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m.bps[gid].exec(gid)
			}
		}
	}
	rplc := m.replication
	for ix := range m.replication.replicas {
		rplc.wg.Go(func() { fn(ix) })
	}
	m.replication.wg.Wait()
}

// Continue - returns to the current state after exiting
func (c *Caller) Continue() { c.continue_sig <- struct{}{} }

// ChangeState - set next executable state to provided one.
// Executable name must be registered through .AddState method
//
// AddState(func(c *ssm.Caller) error {}, "last_state")
//
// where "last_state" is a state name.
func (c *Caller) ChangeState(state string) { c.change_sig <- state }

// ThreadID - returns thread index from multithread mode
func (c *Caller) ThreadID() int { return c.thread_id }

func new_caller(goroutine_id int) *Caller {
	return &Caller{continue_sig: make(chan struct{}, 1), change_sig: make(chan string, 1), thread_id: goroutine_id}
}

func (b *boiler_plate) exec(gid int) {
	b.loop_sleep_sync()
	b.start_h_init(gid)

	clr := new_caller(gid)
	err := b.states[b.last_state](clr)

	if err != nil {
		if eh := b.cfg.Err_handler; eh != nil {
			eh(err)
		}
		return
	}

	select {
	case <-clr.continue_sig:
		return
	case t := <-clr.change_sig:
		if i, ok := b.state_map[t]; ok {
			b.last_state = i
		} else {
			println("[SSM] state not found:", t)
		}
		return
	default:
		b.next_step()
	}
}

func (b *boiler_plate) start_h_init(gid int) {
	start_h_arg := &StartArg{state_name: b.names[b.last_state], thread_id: gid}
	if _sh := b.cfg.Start_handler; _sh != nil {
		_sh(start_h_arg)
	}
}

func (b *boiler_plate) loop_sleep_sync() {
	last_loop := b.last_state+1 == len(b.states)
	if loop_to := b.cfg.Loop_tm; loop_to == 0 && last_loop {
		time.Sleep(loop_to)
	}
}

func (b *boiler_plate) next_step() {
	if b.last_state+1 == len(b.states) {
		b.last_state = 0
		return
	}
	b.last_state += 1
}

func (b *boiler_plate) init_replicated_plates(num int) []*boiler_plate {
	all := make([]*boiler_plate, 0, num)
	for range num {
		_copy := *b
		_copy.states = append([]State(nil), b.states...)
		_copy.names = append([]string(nil), b.names...)
		if b.cfg != nil {
			cfgCopy := *b.cfg
			_copy.cfg = &cfgCopy
		} else {
			_copy.cfg = new(Config)
		}
		_copy.state_map = make(map[string]int)
		for k, v := range b.state_map {
			_copy.state_map[k] = v
		}
		all = append(all, &_copy)
	}
	return all
}

func (s *StartArg) StateName() string { return s.state_name }
func (s *StartArg) ThreadID() int     { return s.thread_id }
