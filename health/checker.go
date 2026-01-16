package health

import (
	"log"
	"net/http"
	"time"

	"github.com/anjy7/Loadbalancer-go/config"
	"github.com/anjy7/Loadbalancer-go/loadbalancer"
)

type HealthChecker struct {
	balancer loadbalancer.LoadBalancer
	backends []config.Backend
	interval time.Duration
}

func NewHealthChecker(balancer loadbalancer.LoadBalancer, backends []config.Backend) *HealthChecker {
	return &HealthChecker{
		balancer: balancer,
		backends: backends,
		interval: 30 * time.Second,
	}
}

func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	go func() {
		for range ticker.C {
			hc.checkHealth()
		}
	}()
}

func (hc *HealthChecker) checkHealth() {
	for _, backend := range hc.backends {
		go func(b config.Backend) {
			hc.checkHTTPHealth(b)
		}(backend)
	}
}

func (hc *HealthChecker) checkHTTPHealth(backend config.Backend) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(backend.URL + "/health")
	if err != nil {
		log.Printf("Backend %s is unhealthy: %v", backend.ID, err)
		hc.balancer.MarkUnhealthy(backend.ID)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		hc.balancer.MarkHealthy(backend.ID)
	} else {
		log.Printf("Backend %s is unhealthy: %v", backend.ID, err)
		hc.balancer.MarkUnhealthy(backend.ID)
	}
}
