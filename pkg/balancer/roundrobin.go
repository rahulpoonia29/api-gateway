package balancer

import (
	"errors"
	"sync"
)

type roundRobinBalancer struct {
	current int
	mutex   sync.Mutex
	targets []string
}

// newRoundrobinBalancer creates a new instance of Roundrobin
func newRoundrobinBalancer(targets []string) *roundRobinBalancer {
	return &roundRobinBalancer{
		current: 0,
		targets: targets,
	}
}

func (rrb *roundRobinBalancer) Elect() (string, error) {
	// Lock the mutex to ensure thread safety
	rrb.mutex.Lock()
	defer rrb.mutex.Unlock()

	if len(rrb.targets) == 0 {
		return "", errors.New("no targets available")
	}

	if len(rrb.targets) == 1 {
		return rrb.targets[0], nil
	}

	println("Current index: ", rrb.current)
	target := rrb.targets[rrb.current]
	rrb.current = (rrb.current + 1) % len(rrb.targets)
	println("Current index after addition: ", rrb.current)

	return target, nil
}
