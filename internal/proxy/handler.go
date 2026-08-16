package proxy

import (
	"encoding/json"
	"errors"
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
	var req CreateProxyRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(
			w,
			"invalid JSON",
			http.StatusBadRequest,
		)
		return
	}

	req.Domain = strings.ToLower(
		strings.TrimSpace(req.Domain),
	)

	req.TargetHost = strings.TrimSpace(
		req.TargetHost,
	)

	if err := validate(req); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	configDir := os.Getenv(
		"NGINX_CONFIG_DIR",
	)

	configPath, err := nginxconfig.CreateOrUpdateProxy(
		configDir,
		nginxconfig.ProxyConfig{
			Domain:     req.Domain,
			TargetHost: req.TargetHost,
			TargetPort: req.TargetPort,
			SSL:        req.SSL,
		},
	)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"failed to apply nginx config: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}

	result := Proxy{
		ID:         req.Domain,
		Domain:     req.Domain,
		TargetHost: req.TargetHost,
		TargetPort: req.TargetPort,
		SSL:        req.SSL,
		ConfigFile: configPath,
	}

	writeJSON(
		w,
		http.StatusCreated,
		result,
	)
}

func ListHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	configDir := os.Getenv(
		"NGINX_CONFIG_DIR",
	)

	configs, err := nginxconfig.ListProxies(
		configDir,
	)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"failed to list nginx proxies: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}

	result := make(
		[]Proxy,
		0,
		len(configs),
	)

	for _, config := range configs {

		result = append(
			result,
			Proxy{
				ID:         config.Domain,
				Domain:     config.Domain,
				TargetHost: config.TargetHost,
				TargetPort: config.TargetPort,
				ConfigFile: nginxconfig.ConfigPath(
					configDir,
					config.Domain,
				),
			},
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		result,
	)
}

func DeleteHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	domain := strings.ToLower(
		strings.TrimSpace(
			r.PathValue("domain"),
		),
	)

	if domain == "" ||
		!isValidHostname(domain) {

		http.Error(
			w,
			"invalid domain",
			http.StatusBadRequest,
		)

		return
	}

	configDir := os.Getenv(
		"NGINX_CONFIG_DIR",
	)

	err := nginxconfig.DeleteProxy(
		configDir,
		domain,
	)

	if errors.Is(
		err,
		nginxconfig.ErrProxyNotFound,
	) {

		http.Error(
			w,
			"proxy not found",
			http.StatusNotFound,
		)

		return
	}

	if err != nil {

		http.Error(
			w,
			fmt.Sprintf(
				"failed to delete proxy: %v",
				err,
			),
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}

func validate(
	req CreateProxyRequest,
) error {

	if req.Domain == "" {
		return &ValidationError{
			"domain is required",
		}
	}

	if !isValidHostname(
		req.Domain,
	) {
		return &ValidationError{
			"domain is invalid",
		}
	}

	if req.TargetHost == "" {
		return &ValidationError{
			"target_host is required",
		}
	}

	if net.ParseIP(req.TargetHost) == nil &&
		!isValidHostname(req.TargetHost) {

		return &ValidationError{
			"target_host must be a valid IP address or hostname",
		}
	}

	if req.TargetPort < 1 ||
		req.TargetPort > 65535 {

		return &ValidationError{
			"target_port must be between 1 and 65535",
		}
	}

	return nil
}

func isValidHostname(
	value string,
) bool {

	if len(value) > 253 {
		return false
	}

	return hostnameRegex.MatchString(
		value,
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(
		data,
	)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}