package main

import (
	"log"
	"net/http"

	"github.com/anjy7/Loadbalancer-go/config"
	"github.com/anjy7/Loadbalancer-go/health"
	"github.com/anjy7/Loadbalancer-go/loadbalancer"
	"github.com/anjy7/Loadbalancer-go/middleware"
	"github.com/anjy7/Loadbalancer-go/proxy"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	var balancer loadbalancer.LoadBalancer
	switch cfg.LoadBalancer.Algorithm {
	case "least_connections":
		balancer = loadbalancer.NewLeastConnectionsBalancer(cfg.Backends)
	default:
		balancer = loadbalancer.NewRoundRobinBalancer(cfg.Backends)
	}

	healthChecker := health.NewHealthChecker(balancer, cfg.Backends)
	healthChecker.Start()

	httpProxy := proxy.NewHTTPProxy(balancer)

	router := mux.NewRouter()
	if cfg.RateLimit.Enable {
		rateLimiter := middleware.NewRateLimiter(
			cfg.RateLimit.RequestsPerSecond,
			cfg.RateLimit.BurstSize,
		)
		router.Use(rateLimiter.Middleware)
	}

	if cfg.JWTSecret != "" {
		authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
		router.Use(authMiddleware.Middleware)
	}

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.PathPrefix("/").Handler(httpProxy)

	log.Printf("Load Balancer starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
