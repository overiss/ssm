package ssm

import (
	"context"
	"time"
)

func CreateMachine() *boiler_plate {
	return &boiler_plate{states: make([]State, 0)}
}

func (b *boiler_plate) WithLoopTimeout(to time.Duration) *boiler_plate {
	b.loop_tm = &to
	return b
}

func (b *boiler_plate) WithErrorHandler(fn func(err error)) *boiler_plate {
	b.err_handler = fn
	return b
}

func (b *boiler_plate) AddState(action State) *boiler_plate {
	b.states = append(b.states, action)
	return b
}

func (b *boiler_plate) UseContext(ctx context.Context) *contextable_bp {
	b.ctx = &contextable{ctx: ctx}
	return &contextable_bp{b}
}

func (c *contextable_bp) Build() *Machine {
	return &Machine{cbp: c}
}

func (m *Machine) Run() {
	for {
		select {
		case <-m.cbp.ctx.ctx.Done():return
		default:m.exec()
		}
	}
}

func (c *Caller) Continue() {c.continue_sig <- struct{}{}}

func new_caller() *Caller {return &Caller{continue_sig: make(chan struct{}, 1)}}

func (m *Machine) exec() {
	m.sleep_sync()
	crucial, clr, exec_chan := m.cbp.states[m.cbp.last_state], new_caller(), make(chan error, 1)
	go func() { defer close(exec_chan); exec_chan <- crucial(clr) }()

	select {
	case e := <-exec_chan:	if _h := m.cbp.err_handler; e != nil && _h != nil {_h(e); return}
	case <-clr.continue_sig:return
	}

	m.next_step()
}

func (m *Machine) sleep_sync() {
	last_loop := m.cbp.last_state+1 == len(m.cbp.states)
	if loop_to := m.cbp.loop_tm; loop_to != nil && last_loop {
		time.Sleep(*loop_to)
	}
}

func (m *Machine) next_step() {
	if m.cbp.last_state+1 == len(m.cbp.states) {m.cbp.last_state = 0;return}
	m.cbp.last_state += 1
}

func (d *data_lake[T]) get() T {
	return d.lake.Get().(T)
}

func (d *data_lake[T]) put(x T) {
	d.lake.Put(x)
}
