// Car example from http://simpy.readthedocs.io/en/latest/simpy_intro/basic_concepts.html
package main

import (
	"fmt"

	"github.com/tebeka/simgo"
)

const (
	parkingDuration = 5
	tripDuration    = 2
)

type Car struct {
	env *simgo.Env
	in  chan int
	out chan int
}

func (c *Car) In() simgo.InChan {
	return c.in
}

func (c *Car) Out() simgo.OutChan {
	return c.out
}

func NewCar(env *simgo.Env) *Car {
	car := &Car{
		env: env,
		// Make channels buffered so env won't block when signaling to start
		in:  make(chan int, 1),
		out: make(chan int, 1),
	}

	go func() {
		for {
			<-car.In()
			fmt.Printf("Start parking at %d\n", env.Now())
			car.Out() <- parkingDuration
			<-car.In()
			fmt.Printf("Start driving at %d\n", env.Now())
			car.Out() <- tripDuration
		}
	}()

	return car
}

func main() {
	env := simgo.NewEnv()
	car := NewCar(env)
	env.Process(car)
	env.Run(15)
}
