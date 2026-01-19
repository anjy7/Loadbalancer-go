package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/anjy7/Loadbalancer-go/loadbalancer"
)

type HTTPProxy struct {
	balancer loadbalancer.LoadBalancer
}

func NewHTTPProxy(balancer loadbalancer.LoadBalancer) *HTTPProxy {
	return &HTTPProxy{
		balancer: balancer,
	}
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := p.balancer.GetNext()
	if backend == nil {
		http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
		return
	}

	target, err := url.Parse(backend.URL)
	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for backend %s: %v", backend.ID, err)
		p.balancer.MarkUnhealthy(backend.ID)
		http.Error(w, "Backend unavailable", http.StatusBadGateway)
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Backend-ID", backend.ID)
	}

	proxy.ServeHTTP(w, r)
}
