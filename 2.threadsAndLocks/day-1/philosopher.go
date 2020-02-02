package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type chopstick struct{ sync.Mutex }

type philosopher struct {
	number    int
	first     *chopstick
	second    *chopstick
	thinkTime int
	eatTime   int
}

func (p *philosopher) run() {
	for {
		// think a while
		fmt.Println("thinking", p.number)
		duration := time.Duration(0)
		if p.thinkTime > 0 {
			duration = time.Duration(rand.Int63n(int64(p.thinkTime)))
		}
		time.Sleep(duration * time.Millisecond)
		p.first.Lock()
		p.second.Lock()
		// eat a while
		fmt.Println("eating", p.number)
		if p.eatTime > 0 {
			duration = time.Duration(rand.Int63n(int64(p.eatTime)))
		}
		time.Sleep(duration * time.Millisecond)
		p.second.Unlock()
		p.first.Unlock()
	}
}

func main() {
	var numOfPhilosophers int
	var think int
	var eat int
	flag.IntVar(&numOfPhilosophers, "p", 1, "num Of Philosophers")
	flag.IntVar(&think, "t", int(rand.Int63n(1000)), "think time")
	flag.IntVar(&eat, "e", int(rand.Int63n(1000)), "eat time")
	flag.Parse()
	numOfChopsticks := numOfPhilosophers
	if numOfPhilosophers == 1 {
		numOfChopsticks = 2
	}

	// init chopsticks
	chopsticks := make([]chopstick, numOfChopsticks)

	philosophers := make([]philosopher, numOfPhilosophers)
	for i := 0; i < numOfPhilosophers; i++ {
		philosopherGetter := func() philosopher {
			if i%2 == 0 {
				return philosopher{
					number:    i,
					first:     &chopsticks[i%numOfChopsticks],
					second:    &chopsticks[(i+1)%numOfChopsticks],
					thinkTime: think,
					eatTime:   eat,
				}
			}
			return philosopher{
				number:    i,
				first:     &chopsticks[(i+1)%numOfChopsticks],
				second:    &chopsticks[i%numOfChopsticks],
				thinkTime: think,
				eatTime:   eat,
			}
		}
		philosophers[i] = philosopherGetter()
		go philosophers[i].run()
	}
	select {}
}
