package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	ssm "github.com/mxmrykov/smart-state-machine"
)

const (
	stateRead  = "read"
	stateWrite = "write"
)

func main() {
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
		}, stateRead).
		AddState(func(c *ssm.Caller) error {
			needsContinue := b(res)
			if needsContinue {
				println("continue")
				c.Continue()
				return nil
			}
			return nil
		}, stateWrite).
		AddState(func(c *ssm.Caller) error {
			r := rand.Intn(10)
			if r > 5 {
				c.ChangeState(stateWrite)
				println("changing state event")
				return nil
			}
			return nil
		}, "last_state").ApplyCfg(&cfg).Build()

	machine.Run()
}

// for future tests
func a() ([]int, error) {
	tm := rand.Intn(10)
	time.Sleep(time.Duration(tm) * time.Second)
	if tm > 5 {
		return nil, errors.New("some error")
	}
	slice := make([]int, 0, tm)
	for i := range tm {
		slice = append(slice, i*2)
	}
	return slice, nil
}

func b(a []int) bool {
	fmt.Println(a)
	return rand.Intn(2)%2 == 0
}
