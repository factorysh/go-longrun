package sse

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvents(t *testing.T) {
	evts := NewEvents()
	wait := &sync.WaitGroup{}
	wait.Add(2)
	for i := 0; i < 2; i++ {
		go func(n int) {
			ctx := context.TODO()
			events := evts.Subscribe(ctx, 0)
			cpt := 0
			for {
				evt := <-events
				fmt.Println(n, "evt : ", evt)
				cpt++
				if cpt == 2 {
					break
				}
			}
			wait.Done()
		}(i)
	}
	evts.Append(&Event{
		Data: "Pim",
	})
	evts.Append(&Event{
		Data: "Pam",
	})
	wait.Wait()
	evts.Close()
	assert.Len(t, evts.broadcast, 0)
}
