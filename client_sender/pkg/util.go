package pkg

import "time"

type Throttler struct {
	delay time.Duration
	timer *time.Timer
}

func NewThrottler(delay time.Duration) *Throttler {
	return &Throttler{delay: delay}
}

func (t *Throttler) Call(f func()) {
	if t.timer != nil {
		t.timer.Stop()
	}
	t.timer = time.AfterFunc(t.delay, f)
}
