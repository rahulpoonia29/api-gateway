package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/rahul/api-gateway/pkg/balancer"
	"github.com/rahul/api-gateway/pkg/config"
	"github.com/rahul/api-gateway/utils"
)

type HTTPHandler struct {
	app *utils.App
}

// ServeHTTP handles incoming HTTP requests.
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	fmt.Printf("Received request: %s %s\n", r.Method, path)

	// 1. Find the service configuration based on the path
	key, value, found := h.app.RouteTree.LongestPrefix(path)
	if !found {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	serviceConfig, ok := value.(config.ServiceConfig)
	if !ok {
		http.Error(w, "Internal gateway error", http.StatusInternalServerError)
		fmt.Printf("Error: Could not cast route value to ServiceConfig\n")
		return
	}

	// 2. **Middleware Chain (Placeholder):**
	// In the future, we'll apply middleware here based on serviceConfig
	// For now, we'll just proceed to proxying.

	// 3. Proxying: Forward the request to the upstream service
	h.proxyRequest(w, r, serviceConfig, key)
}

// proxyRequest forwards the request to the upstream service.
func (h *HTTPHandler) proxyRequest(w http.ResponseWriter, r *http.Request, serviceConfig config.ServiceConfig, prefixPath string) {
	if len(serviceConfig.Proxy.Upstream.Targets) == 0 {
		http.Error(w, "Service has no upstream targets", http.StatusInternalServerError)
		return
	}

	// Implement load balancing strategy
	balancer, err := balancer.NewBalancer(&serviceConfig.Proxy.Upstream)
	if err != nil {
		http.Error(w, "Error creating balancer", http.StatusInternalServerError)
		return
	}

	targetURL, err := balancer.Elect(serviceConfig.Proxy.Upstream.Targets)
	if err != nil {
		http.Error(w, "Error selecting target", http.StatusInternalServerError)
		return
	}
	targetURL = strings.TrimSuffix(targetURL, "/")

	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid upstream URL", http.StatusInternalServerError)
		fmt.Printf("Error parsing target URL '%s': %v\n", targetURL, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Prepare to modify the request *before* forwarding
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req) // Call the original director

		// Modify request path based on configuration (example: strip prefix)
		remainingPath := strings.TrimPrefix(req.URL.Path, prefixPath)
		if !strings.HasPrefix(remainingPath, "/") && remainingPath != "" {
			remainingPath = "/" + remainingPath
		}
		req.URL.Path = remainingPath

		// Optionally, you could modify headers here if needed based on serviceConfig
		// req.Header.Set("X-Gateway-Version", "1.0") // Example
	}

	// Log forwarding action
	fmt.Printf("Forwarding to: %s%s\n", targetURL, r.URL.Path)

	proxy.ServeHTTP(w, r)
}
