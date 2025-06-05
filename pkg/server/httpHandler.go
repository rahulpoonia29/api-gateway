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
	app       *utils.App
	balancers map[string]balancer.Balancer
}

func NewHTTPHandler(app *utils.App) *HTTPHandler {
	return &HTTPHandler{
		app:       app,
		balancers: make(map[string]balancer.Balancer, 0),
	}
}

// ServeHTTP handles incoming HTTP requests.
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	h.app.Logger.Debug("received request", "method", r.Method, "path", path)

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
	existingBalancer, ok := h.balancers[serviceConfig.Name]
	if !ok {
		newBalancer, err := balancer.NewBalancer(&serviceConfig.Proxy.Upstream)
		if err != nil {
			h.app.Logger.Error("error creating balancer",
				"service", serviceConfig.Name,
				"strategy", serviceConfig.Proxy.Upstream.Balancing,
				"error", err)
			http.Error(w, "Error creating balancer", http.StatusInternalServerError)
			return
		}
		h.balancers[serviceConfig.Name] = newBalancer
		existingBalancer = newBalancer
	}

	targetURL, err := existingBalancer.Elect()
	if err != nil {
		h.app.Logger.Error("error selecting target",
			"service", serviceConfig.Name,
			"error", err)
		http.Error(w, "Error selecting upstream target", http.StatusInternalServerError)
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

	// Modify the request path based on configuration
	proxy.Director = func(req *http.Request) {
		if !strings.HasSuffix(prefixPath, "/") {
			prefixPath += "/"
		}

		remainingPath := strings.TrimPrefix(req.URL.Path, prefixPath)

		if !strings.HasSuffix(remainingPath, "/") {
			remainingPath += "/"
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = remainingPath

		if target.RawQuery != "" {
			req.URL.RawQuery = target.RawQuery
		}

		// Forward all original headers
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Set custom gateway headers
		req.Header.Set("X-Gateway-Version", "0.1")
		req.Header.Set("X-Gateway-Service", serviceConfig.Name)

	}

	// Log forwarding action
	h.app.Logger.Info("forwarding request",
		"service", serviceConfig.Name,
		"target", targetURL,
		"path", r.URL.Path)

	// Use the reverse proxy to forward the request to the upstream target
	proxy.ServeHTTP(w, r)
}
