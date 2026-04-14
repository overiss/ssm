package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	ssm "github.com/mxmrykov/smart-state-machine"
)

func main() {
	ctx := context.Background()

	machine := ssm.CreateMachine().
		AddState(func(c *ssm.Caller) error {
			fmt.Println("starting state 1")
			r := rand.Intn(10)
			if r > 5 {
				return errors.New("some error")
			}
			fmt.Println("ending state 1")
			return nil
		}).
		AddState(func(c *ssm.Caller) error {
			fmt.Println("starting state 2")
			r := rand.Intn(10)
			if r > 5 {
				return errors.New("some another error")
			}
			fmt.Println("ending state 2")
			return nil
		}).
		AddState(func(c *ssm.Caller) error {
			fmt.Println("starting state 3")
			r := rand.Intn(10)
			if r > 5 {
				return errors.New("some moreover error")
			}
			fmt.Println("ending state 3")
			return nil
		}).
		WithLoopTimeout(time.Second).
		WithErrorHandler(func(err error) {
			fmt.Printf("ERROR! %v\n", err)
		}).
		UseContext(ctx).Build()

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

func b(a []int) {
	time.Sleep(time.Second)
	fmt.Println(a)
}
