package balancer

import (
	"errors"
	"sync"
)

type roundRobinBalancer struct {
	current int
	mutex   sync.Mutex
}

// newRoundrobinBalancer creates a new instance of Roundrobin
func newRoundrobinBalancer() *roundRobinBalancer {
	return &roundRobinBalancer{}
}

func (rrb *roundRobinBalancer) Elect(targets []string) (string, error) {
	// Lock the mutex to ensure thread safety
	rrb.mutex.Lock()
	defer rrb.mutex.Unlock()

	if len(targets) == 0 {
		return "", errors.New("no targets available")
	}

	if len(targets) == 1 {
		return targets[0], nil
	}

	if rrb.current >= len(targets) {
		rrb.current = 0
	}

	target := targets[rrb.current]
	rrb.current++

	return target, nil
}
