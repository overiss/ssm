# Smart State Machine

A wrapper for the state machine that allows you to scale all necessary activities and move them to different service layers.

## Usage example

```go
ctx := context.Background()
cfg := ssm.Config{
	Loop_tm: time.Second,
	Err_handler: func(err error) {
		fmt.Printf("ERROR! %v\n", err)
	},
	Start_handler: func(s *ssm.StartArg) {
		fmt.Printf("Starting state %s at thread %d\n", s.StateName(), s.ThreadID())
	},
	Threads: 2,
}
var (
	res []int
	err error
)

machine := ssm.CreateMachine().
	AddState(func(c *ssm.Caller) error {
		res, err = a()
		if err != nil {
			return err
		}
		return nil
	}, "state_read").
	AddState(func(c *ssm.Caller) error {
		needsContinue := b(res)
		if needsContinue {
			println("threadID: ", c.ThreadID(), ", continue")
			c.Continue()
			return nil
		}
		return nil
	}, "state_write").
	AddState(func(c *ssm.Caller) error {
		r := rand.Intn(10)
		if r > 5 {
			c.ChangeState(stateWrite)
			println("threadID: ", c.ThreadID(), ", changing state event")
			return nil
		}
		return nil
	}, "last_state").ApplyCfg(&cfg).Build()

machine.Run(ctx)
```
