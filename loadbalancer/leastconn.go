package loadbalancer

import (
	"math"
	"sync"

	"github.com/anjy7/Loadbalancer-go/config"
)

type LeastConnectionsBalancer struct {
	backends    []config.Backend
	connections map[string]int
	mutex       sync.RWMutex
}

func NewLeastConnectionsBalancer(backends []config.Backend) *LeastConnectionsBalancer {
	connections := make(map[string]int)
	for _, backend := range backends {
		connections[backend.ID] = 0
	}

	return &LeastConnectionsBalancer{
		backends:    backends,
		connections: connections,
	}
}

func (lb *LeastConnectionsBalancer) GetNext() *config.Backend {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	var selected *config.Backend
	minConnections := math.MaxInt

	for i := range lb.backends {
		if !lb.backends[i].Health {
			continue
		}

		connections := lb.connections[lb.backends[i].ID]
		if connections < minConnections {
			minConnections = connections
			selected = &lb.backends[i]
		}
	}

	if selected != nil {
		lb.connections[selected.ID]++
	}

	return selected
}

func (lb *LeastConnectionsBalancer) DecrementConnections(backendID string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if count, exists := lb.connections[backendID]; exists && count > 0 {
		lb.connections[backendID]--
	}
}

func (lb *LeastConnectionsBalancer) MarkUnhealthy(backendID string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i := range lb.backends {
		if lb.backends[i].ID == backendID {
			lb.backends[i].Health = false
			break
		}
	}
}

func (lb *LeastConnectionsBalancer) MarkHealthy(backendID string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i := range lb.backends {
		if lb.backends[i].ID == backendID {
			lb.backends[i].Health = true
			break
		}
	}
}

func (lb *LeastConnectionsBalancer) UpdateBackends(backends []config.Backend) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.backends = backends
}
