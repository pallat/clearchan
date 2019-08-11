package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	w := worker{
		chsignal: make(chan struct{}),
	}

	o := observer{
		chsignal: make(chan struct{}),
		wg:       sync.WaitGroup{},
	}

	go w.dosomething()
	go o.alert()

	go func() {
		for {
			o.listen() <- <-w.signal()
		}
	}()

	o.wait(5)
}

type worker struct {
	chsignal chan struct{}
}

func (w *worker) dosomething() {
	rand.Seed(time.Now().UnixNano())
	c := time.Tick(time.Millisecond * 10)
	for range c {
		if (rand.Int() % 10) == 0 {
			w.chsignal <- struct{}{}
		}
	}
}
func (w *worker) signal() <-chan struct{} {
	return w.chsignal
}

type observer struct {
	chsignal chan struct{}
	wg       sync.WaitGroup
}

func (o *observer) alert() {
	for range o.chsignal {
		fmt.Println("got signal")
		o.wg.Done()
	}
}

func (o *observer) listen() chan<- struct{} {
	return o.chsignal
}

func (o *observer) wait(n int) {
	o.wg.Add(n)
	o.wg.Wait()
}
