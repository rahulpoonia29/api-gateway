package server

import (
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

	h.app.Logger.Debug("received request",
		"method", r.Method,
		"path", path,
		"remoteAddr", r.RemoteAddr)

	// 1. Find the service configuration based on the path
	key, value, found := h.app.RouteTree.LongestPrefix(path)
	if !found {
		h.app.Logger.Warn("service not found for path", "path", path)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	serviceConfig, ok := value.(config.ServiceConfig)
	if !ok {
		h.app.Logger.Error("failed to cast route value to ServiceConfig", "path", path)
		http.Error(w, "Internal gateway error", http.StatusInternalServerError)
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
		h.app.Logger.Error("service has no upstream targets", "service", serviceConfig.Name)
		http.Error(w, "Service has no upstream targets", http.StatusInternalServerError)
		return
	}

	// Implement load balancing strategy
	balancer, err := balancer.NewBalancer(&serviceConfig.Proxy.Upstream)
	if err != nil {
		h.app.Logger.Error("error creating balancer",
			"service", serviceConfig.Name,
			"strategy", serviceConfig.Proxy.Upstream.Balancing,
			"error", err)
		http.Error(w, "Error creating balancer", http.StatusInternalServerError)
		return
	}

	targetURL, err := balancer.Elect(serviceConfig.Proxy.Upstream.Targets)
	if err != nil {
		h.app.Logger.Error("error selecting target",
			"service", serviceConfig.Name,
			"error", err)
		http.Error(w, "Error selecting target", http.StatusInternalServerError)
		return
	}
	targetURL = strings.TrimSuffix(targetURL, "/")

	target, err := url.Parse(targetURL)
	if err != nil {
		h.app.Logger.Error("invalid upstream URL",
			"service", serviceConfig.Name,
			"url", targetURL,
			"error", err)
		http.Error(w, "Invalid upstream URL", http.StatusInternalServerError)
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
	h.app.Logger.Info("forwarding request",
		"service", serviceConfig.Name,
		"target", targetURL,
		"path", r.URL.Path)

	proxy.ServeHTTP(w, r)
}
