package status

import (
	"encoding/json"
	"net/http"
	"os"

	nginxmanager "nginx-manager-service/internal/nginx"
)

type Response struct {
	Service ServiceStatus     `json:"service"`
	Nginx   NginxStatus       `json:"nginx"`
	Cert    CertificateStatus `json:"certificate"`
	Proxies ProxyStatus       `json:"proxies"`
}

type ServiceStatus struct {
	Status string `json:"status"`
}

type NginxStatus struct {
	ConfigValid bool   `json:"config_valid"`
	ReloadMode  string `json:"reload_mode"`
}

type CertificateStatus struct {
	Mode string `json:"mode"`
}

type ProxyStatus struct {
	Managed int `json:"managed"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	configDir := os.Getenv("NGINX_CONFIG_DIR")

	/*
		Check nginx config.
	*/

	configValid := true

	if err := nginxmanager.TestConfig(); err != nil {
		configValid = false
	}

	/*
		Get reload mode.
	*/

	reloadMode := os.Getenv("NGINX_RELOAD_MODE")

	if reloadMode == "" {
		reloadMode = "noop"
	}

	/*
		Get certificate mode.
	*/

	certificateMode := os.Getenv("CERTIFICATE_MODE")

	if certificateMode == "" {
		certificateMode = "disabled"
	}

	/*
		Count managed proxies.
	*/

	managedProxies := 0

	proxies, err := nginxmanager.ListProxies(configDir)

	if err == nil {
		managedProxies = len(proxies)
	}

	response := Response{
		Service: ServiceStatus{
			Status: "online",
		},

		Nginx: NginxStatus{
			ConfigValid: configValid,
			ReloadMode:  reloadMode,
		},

		Cert: CertificateStatus{
			Mode: certificateMode,
		},

		Proxies: ProxyStatus{
			Managed: managedProxies,
		},
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}