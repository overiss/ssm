package ssm

import (
	"context"
	"fmt"
	"time"
)

// CreateMachine - creates a set of methods for configuring the state machine
func CreateMachine(ctx context.Context) *boiler_plate {
	return &boiler_plate{states: make([]State, 0), cfg: new(Config), ctx: ctx}
}

// ApplyCfg - applies configurations to the state machine
func (b *boiler_plate) ApplyCfg(cfg *Config) *boiler_plate {b.cfg = cfg;return b}

// AddState - adds an executable state to the machine. You can specify a custom
// state name as the second argument, otherwise it will be determined by default.
func (b *boiler_plate) AddState(s State, custom_name ...string) *boiler_plate {
	b.states = append(b.states, s)
	name := fmt.Sprintf(stateNameMask, len(b.states))
	if len(custom_name) > 0 {name = custom_name[0]}
	b.names = append(b.names, name)
	return b
}

// Build - assembles the configured machine into an executable method
func (b *boiler_plate) Build() *Machine {return &Machine{bp: b}}

func (m *Machine) Run() {
	for {
		select {
		case <-m.bp.ctx.Done():return
		default:m.exec()
		}
	}
}

// Continue - returns to the current state after exiting
func (c *Caller) Continue() {c.continue_sig <- struct{}{}}

func new_caller() *Caller {return &Caller{continue_sig: make(chan struct{}, 1)}}

func (m *Machine) exec() {
	m.loop_sleep_sync()
	if _sh := m.bp.cfg.Start_handler; _sh != nil {_sh(m.bp.names[m.bp.last_state])}
	crucial, clr, exec_chan := m.bp.states[m.bp.last_state], new_caller(), make(chan error, 1)
	go func() { defer close(exec_chan); exec_chan <- crucial(clr) }()

	select {
	case e := <-exec_chan:	if _eh := m.bp.cfg.Err_handler; e != nil && _eh != nil {_eh(e); return}
	case <-clr.continue_sig:return
	}

	m.next_step()
}

func (m *Machine) loop_sleep_sync() {
	last_loop := m.bp.last_state+1 == len(m.bp.states)
	if loop_to := m.bp.cfg.Loop_tm; loop_to == 0 && last_loop {time.Sleep(loop_to)}
}

func (m *Machine) next_step() {
	if m.bp.last_state+1 == len(m.bp.states) {m.bp.last_state = 0;return}
	m.bp.last_state += 1
}
