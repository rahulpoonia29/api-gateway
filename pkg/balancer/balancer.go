package balancer

import (
	"errors"

	"github.com/rahul/api-gateway/pkg/config"
)

type Balancer interface {
	Elect() (string, error)
}

func NewBalancer(upstream *config.UpstreamConfig) (Balancer, error) {
	var balancer Balancer
	var err error = nil

	switch upstream.Balancing {
	case config.RoundRobin:
		balancer = newRoundrobinBalancer(upstream.Targets)
	default:
		err = errors.New("unsupported load balancing strategy")

	}

	return balancer, err
}
