package example

import (
	"context"
	"time"
)

func waitReady(ready func() bool) {
	for !ready() {
		// ruleid: no-sleep-poll-loop
		time.Sleep(50 * time.Millisecond)
	}
}

func waitReadyCtx(ctx context.Context, ready func() bool) error {
	t := time.NewTicker(50 * time.Millisecond)
	defer t.Stop()
	for !ready() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		// ok: no-sleep-poll-loop
		case <-t.C:
		}
	}
	return nil
}

func waitSocket(exists func() bool) bool {
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if exists() {
			return true
		}
		// ok: no-sleep-poll-loop
		time.Sleep(50 * time.Millisecond)
	}
	return false
}
