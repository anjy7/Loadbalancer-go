package loadbalancer

import (
	"github.com/anjy7/Loadbalancer-go/config"
)

type LoadBalancer interface {
	GetNext() *config.Backend
	UpdateBackends(backends []config.Backend)
	MarkUnhealthy(backendID string)
	MarkHealthy(backendID string)
}
