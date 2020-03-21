package main

import (
	"fmt"
	"reflect"
	"sync"
)

func main() {
	ch1 := asChan(1, 2, 3, 4, 5, 6, 7, 8)
	ch2 := asChan(11, 12, 13, 14, 15, 16, 17, 18)
	ch3 := asChan(21, 22, 23, 24, 25, 26, 27, 28)
	for val := range mergeReflect(ch1, ch2, ch3) {
		fmt.Println(val)
	}
}

func merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))
		for _, c := range chans {
			c := c
			go func() {
				for v := range c {
					out <- v
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func mergeReflect(chans ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface().(int)
		}
	}()
	return out
}

func mergeRec(chans ...<-chan int) <-chan int {
	n := len(chans)
	switch n {
	case 0:
		return nil
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		return mergeTwo(merge(chans[:n/2]...), merge(chans[n/2:]...))
	}
}

func mergeTwo(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func asChan(vs ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range vs {
			out <- v
		}
		close(out)
	}()
	return out
}
