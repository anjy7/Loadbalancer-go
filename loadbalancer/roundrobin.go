package loadbalancer

import (
	"sync"

	"github.com/anjy7/Loadbalancer-go/config"
)

type RoundRobinBalancer struct {
	backends []config.Backend
	current  uint64
	mutex    sync.RWMutex
}

func NewRoundRobinBalancer(backends []config.Backend) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		backends: backends,
		current:  0,
	}
}

func (rb *RoundRobinBalancer) GetNext() *config.Backend {
	rb.mutex.RLock()
	defer rb.mutex.RUnlock()

	if len(rb.backends) == 0 {
		return nil
	}
	current := rb.current
	current++
	for {
		if current == rb.current {
			if rb.backends[current].Health {
				return &rb.backends[current]
			} else {
				return nil
			}
		}
		if rb.backends[current].Health {
			return &rb.backends[current]
		}

		current++
		if current > uint64(len(rb.backends)) {
			current = 0
		}

	}
}

func (rb *RoundRobinBalancer) UpdateBackends(backends []config.Backend) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()
	rb.backends = backends
}

func (rb *RoundRobinBalancer) MarkUnhealthy(backendID string) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	for i := range rb.backends {
		if rb.backends[i].ID == backendID {
			rb.backends[i].Health = false
			break
		}
	}
}

func (rb *RoundRobinBalancer) MarkHealthy(backendID string) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	for i := range rb.backends {
		if rb.backends[i].ID == backendID {
			rb.backends[i].Health = true
			break
		}
	}
}
