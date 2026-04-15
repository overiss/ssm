# Smart State Machine

A wrapper for the state machine that allows you to scale all necessary activities and move them to different service layers.

## Usage example

```
ctx := context.Background()
cfg := ssm.Config{
	Loop_tm: time.Second,
	Err_handler: func(err error) {
		fmt.Printf("ERROR! %v\n", err)
	},
	Start_handler: func(state_name string) {
		fmt.Printf("Starting state %s\n", state_name)
	},
}
var (
	res []int
	err error
)

machine := ssm.CreateMachine(ctx).
	AddState(func(c *ssm.Caller) error {
		res, err = a()
		if err != nil {
			return err
		}
		return nil
	}).
	AddState(func(c *ssm.Caller) error {
		needsContinue := b(res)
		if needsContinue {
			fmt.Println("continue")
			c.Continue()
			return nil
		}
		fmt.Println("aft")
		return nil
	}).
	AddState(func(c *ssm.Caller) error {
		r := rand.Intn(10)
		if r > 5 {
			return errors.New("some moreover error")
		}
		return nil
	}).ApplyCfg(&cfg).Build()

machine.Run()
```
