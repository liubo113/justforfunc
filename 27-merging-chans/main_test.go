package main

import (
	"fmt"
	"testing"
)

var funcs = []struct {
	name string
	f    func(...<-chan int) <-chan int
}{
	{"goroutines", merge},
	{"reflection", mergeReflect},
	{"recursion", mergeRec},
}

func TestMerge(t *testing.T) {
	for _, f := range funcs {
		t.Run(f.name, func(t *testing.T) {
			c := f.f(asChan(1, 2, 3), asChan(4, 5, 6), asChan(7, 8, 9))
			seen := make(map[int]struct{})
			for v := range c {
				if _, ok := seen[v]; ok {
					t.Errorf("duplicated val %d", v)
					continue
				}
				seen[v] = struct{}{}
			}
			for i := 1; i <= 9; i++ {
				if _, ok := seen[i]; !ok {
					t.Errorf("did not found val %d", i)
				}
			}
		})
	}
}

func BenchmarkMerge(b *testing.B) {
	for _, f := range funcs {
		for n := 1; n <= 1024; n *= 2 {
			chans := make([]<-chan int, n)
			b.Run(fmt.Sprintf("%s/%d", f.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for i := range chans {
						chans[i] = asChan(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
					}
					b.StartTimer()
					for range f.f(chans...) {
					}
					b.StopTimer()
				}
			})
		}
	}
}
