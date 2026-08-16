package nginx

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var serverNameRegex = regexp.MustCompile(
	`(?m)^\s*server_name\s+([^\s;]+)\s*;`,
)

var proxyPassRegex = regexp.MustCompile(
	`(?m)^\s*proxy_pass\s+http://([^;]+)\s*;`,
)

func ListProxies(configDir string) ([]ProxyConfig, error) {
	if configDir == "" {
		return nil, fmt.Errorf(
			"nginx config directory is not configured",
		)
	}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf(
			"read nginx config directory: %w",
			err,
		)
	}

	proxies := make([]ProxyConfig, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".conf" {
			continue
		}

		if !strings.HasPrefix(
			entry.Name(),
			"nginx-manager--",
		) {
			continue
		}

		path := filepath.Join(
			configDir,
			entry.Name(),
		)

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf(
				"read config %s: %w",
				entry.Name(),
				err,
			)
		}

		if !bytes.Contains(
			content,
			[]byte("# managed-by: nginx-manager-service"),
		) {
			continue
		}

		config, err := parseProxyConfig(content)
		if err != nil {
			continue
		}

		proxies = append(
			proxies,
			config,
		)
	}

	return proxies, nil
}

func parseProxyConfig(content []byte) (ProxyConfig, error) {
	serverMatch := serverNameRegex.FindSubmatch(content)
	if len(serverMatch) < 2 {
		return ProxyConfig{}, fmt.Errorf(
			"server_name not found",
		)
	}

	proxyMatch := proxyPassRegex.FindSubmatch(content)
	if len(proxyMatch) < 2 {
		return ProxyConfig{}, fmt.Errorf(
			"proxy_pass not found",
		)
	}

	upstream := strings.TrimSpace(
		string(proxyMatch[1]),
	)

	host, portString, err := net.SplitHostPort(upstream)
	if err != nil {
		return ProxyConfig{}, fmt.Errorf(
			"invalid proxy_pass upstream: %w",
			err,
		)
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return ProxyConfig{}, fmt.Errorf(
			"invalid target port: %w",
			err,
		)
	}

	return ProxyConfig{
		Domain:     string(serverMatch[1]),
		TargetHost: host,
		TargetPort: port,
	}, nil
}