package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port         string             `json:"port"`
	JWTSecret    string             `json:"jwt_secret"`
	LoadBalancer LoadBalancerConfig `json:"load_balancer"`
	RateLimit    RateLimitConfig    `json:"rate_limit"`
	Backends     []Backend          `json:"backends"`
}

type LoadBalancerConfig struct {
	Algorithm string `json:"algorithm"` // "round_robin" or "least_connections"
}

type RateLimitConfig struct {
	Enable            bool `json:"enable"`
	RequestsPerSecond int  `json:"requests_per_second"`
	BurstSize         int  `json:"burst_size"`
}

type Backend struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Health bool   `json:"health"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)

	return config, err
}
