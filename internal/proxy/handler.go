package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	nginxconfig "nginx-manager-service/internal/nginx"
)

var hostnameRegex = regexp.MustCompile(
	`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`,
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	req.Domain = strings.ToLower(strings.TrimSpace(req.Domain))
	req.TargetHost = strings.TrimSpace(req.TargetHost)

	if err := validate(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	configDir := os.Getenv("NGINX_CONFIG_DIR")

	configPath, err := nginxconfig.CreateOrUpdateProxy(
		configDir,
		nginxconfig.ProxyConfig{
			Domain:     req.Domain,
			TargetHost: req.TargetHost,
			TargetPort: req.TargetPort,
		},
	)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to generate nginx config: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	proxy := Proxy{
		ID:         req.Domain,
		Domain:     req.Domain,
		TargetHost: req.TargetHost,
		TargetPort: req.TargetPort,
		ConfigFile: configPath,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(proxy)
}

func validate(req CreateProxyRequest) error {
	if req.Domain == "" {
		return &ValidationError{"domain is required"}
	}

	if !isValidHostname(req.Domain) {
		return &ValidationError{"domain is invalid"}
	}

	if req.TargetHost == "" {
		return &ValidationError{"target_host is required"}
	}

	if net.ParseIP(req.TargetHost) == nil && !isValidHostname(req.TargetHost) {
		return &ValidationError{
			"target_host must be a valid IP address or hostname",
		}
	}

	if req.TargetPort < 1 || req.TargetPort > 65535 {
		return &ValidationError{
			"target_port must be between 1 and 65535",
		}
	}

	return nil
}

func isValidHostname(value string) bool {
	if len(value) > 253 {
		return false
	}

	return hostnameRegex.MatchString(value)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}