package balancer

import (
	"errors"

	"github.com/rahul/api-gateway/pkg/config"
)

type Balancer interface {
	Elect(targets []string) (string, error)
}

func NewBalancer(upstream *config.UpstreamConfig) (Balancer, error) {
	var balancer Balancer
	var err error = nil

	switch upstream.Balancing {
	case config.RoundRobin:
		balancer = newRoundrobinBalancer()
	default:
		err = errors.New("unsupported load balancing strategy")

	}

	return balancer, err
}
