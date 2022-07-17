package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type SMSP struct {
	sync.WaitGroup
	c chan int
}

func NewSMSP(chanSize int) *SMSP {
	ret := &SMSP{
		c: make(chan int, chanSize),
	}

	return ret
}

func (s *SMSP) StartProduce(ctx context.Context) {
	s.Add(1)
	defer s.Done()
	ticker := time.NewTicker(1 * time.Second)
	rand.Seed(time.Now().Unix())
	for {
		select {
		case <-ticker.C:
			v := rand.Intn(10)
			fmt.Printf("Produce %d\n", v)
			s.c <- v
		case <-ctx.Done():
			close(s.c)
			return
		}
	}
}

func (s *SMSP) StartConsume(id string) {
	s.Add(1)
	defer s.Done()
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		v, ok := <-s.c
		if !ok {
			return
		}
		fmt.Printf("Consume %s %d\n", id, v)
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
		select {
		case sig := <-signalChan:
			fmt.Printf("Receive signal '%s', exit\n", sig.String())

			cancel()
		}
	}()

	s := NewSMSP(10)
	go s.StartConsume("consumer1")
	s.StartProduce(ctx)
	s.Wait()
}
