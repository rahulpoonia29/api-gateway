package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/rahul/api-gateway/pkg/config"
	"github.com/rahul/api-gateway/utils"
)

type HTTPHandler struct {
	app *utils.App
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	fmt.Printf("Received request: %s %s\n", r.Method, path)

	// Find the longest matching prefix in our route tree
	key, value, found := h.app.RouteTree.LongestPrefix(path)

	if !found {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Extract the service config from the value
	service, ok := value.(config.ServiceConfig)
	if !ok {
		http.Error(w, "Internal gateway error", http.StatusInternalServerError)
		fmt.Printf("Error: Could not cast route value to ServiceConfig\n")
		return
	}

	// Strip the prefix to get the remaining path
	remainingPath := strings.TrimPrefix(path, key)
	// Add a leading slash if the remaining path is not empty and does not start with a slash
	if !strings.HasPrefix(remainingPath, "/") && remainingPath != "" {
		remainingPath = "/" + remainingPath
	}

	if len(service.Proxy.Upstream.Targets) == 0 {
		http.Error(w, "Service has no upstream targets", http.StatusInternalServerError)
		return
	}

	targetURL := service.Proxy.Upstream.Targets[0]
	targetURL = strings.TrimSuffix(targetURL, "/")

	// Forward the request to the upstream service
	forwardRequest(w, r, targetURL, remainingPath)
}

// Forward the request to the upstream target
func forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string, remainingPath string) {
	// Parse the target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid upstream URL", http.StatusInternalServerError)
		fmt.Printf("Error parsing target URL '%s': %v\n", targetURL, err)
		return
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Update the request URL
	r.URL.Path = remainingPath
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme

	// Update the request Host header to match the target
	r.Host = target.Host

	// Log the forwarding
	fmt.Printf("Forwarding to: %s%s\n", targetURL, remainingPath)

	// Forward the request
	proxy.ServeHTTP(w, r)
}
