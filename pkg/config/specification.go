package config

// GatewayConfig represents the main configuration structure for the API gateway
type GatewayConfig struct {
	Gateway  GatewaySettings `json:"gateway"`
	Services []ServiceConfig `json:"services"`
}

// GatewaySettings contains basic settings for the gateway
type GatewaySettings struct {
	Port int `json:"port"`
}

// ServiceConfig represents configuration for an API service
type ServiceConfig struct {
	Name        string      `json:"name"`
	Active      bool        `json:"active"`
	Description string      `json:"description,omitempty"`
	Proxy       ProxyConfig `json:"proxy"`
	// Authentication *AuthenticationConfig `json:"authentication,omitempty"`
	// RateLimit      *RateLimitConfig      `json:"rateLimit,omitempty"`
	// HealthCheck    *HealthCheckConfig    `json:"healthCheck,omitempty"`
}

// ProxyConfig defines how requests are proxied to backend services
type ProxyConfig struct {
	ListenPath string         `json:"listenPath"`
	Upstream   UpstreamConfig `json:"upstream"`
	StripPath  bool           `json:"stripPath"`
	AppendPath bool           `json:"appendPath"`
	Methods    []string       `json:"methods"` // HttpMethod values: GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD
}

// UpstreamConfig defines upstream service configuration
type UpstreamConfig struct {
	Balancing BalancingStrategy `json:"balancing"` // Load balancing algorithm to use
	Targets   []string          `json:"targets"`
}

// BalancingStrategy defines the available load balancing algorithms
type BalancingStrategy string

const (
	RoundRobin BalancingStrategy = "roundrobin"
	LeastConn  BalancingStrategy = "least_conn"
	IPHash     BalancingStrategy = "ip_hash"
)

// AuthenticationConfig defines authentication settings for a route
// type AuthenticationConfig struct {
// 	Plugin string                 `json:"plugin"`
// 	Config map[string]interface{} `json:"config"`
// }

// RateLimitConfig defines rate limiting settings for a route
// type RateLimitConfig struct {
// 	Identifier RateLimitIdentifierConfig `json:"identifier"`
// 	Plugin     string                    `json:"plugin"`
// 	Config     map[string]interface{}    `json:"config"`
// }

// RateLimitIdentifierConfig defines how to identify clients for rate limiting
// type RateLimitIdentifierConfig struct {
// 	Source string `json:"source"` // ip, header, jwtClaim, apiKey
// 	Key    string `json:"key,omitempty"`
// }

// HealthCheckConfig defines health check settings for a route
// type HealthCheckConfig struct {
// 	URL     string `json:"url"`
// 	Timeout int    `json:"timeout"`
// }
