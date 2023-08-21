package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var instance *Closer = nil

type Closer struct {
	mu    sync.Mutex
	funcs []func(ctx context.Context) error
}

func (c *Closer) Add(f func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		msgs     = make([]string, 0, len(c.funcs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				msgs = append(msgs, fmt.Sprintf("[!] %v", err))
			}
		}
		complete <- struct{}{}
	}()
	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(msgs, "\n"),
		)
	}
	return nil
}

func GetInstance() *Closer {
	var once sync.Once
	once.Do(func() {
		instance = new(Closer)
	})
	return instance
}
