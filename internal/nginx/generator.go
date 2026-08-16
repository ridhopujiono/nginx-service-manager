package nginx

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type ProxyConfig struct {
	Domain      string
	TargetHost  string
	TargetPort  int
	SSL         bool
	ACMEWebroot string
}

const httpProxyTemplate = `# managed-by: nginx-manager-service

server {
    listen 80;
    server_name {{.Domain}};

    location ^~ /.well-known/acme-challenge/ {
        root {{.ACMEWebroot}};
        default_type text/plain;
        try_files $uri =404;
    }

    location / {
        proxy_pass http://{{.TargetHost}}:{{.TargetPort}};

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
`

const httpsProxyTemplate = `# managed-by: nginx-manager-service

server {
    listen 80;
    server_name {{.Domain}};

    location ^~ /.well-known/acme-challenge/ {
        root {{.ACMEWebroot}};
        default_type text/plain;
        try_files $uri =404;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name {{.Domain}};

    ssl_certificate /etc/letsencrypt/live/{{.Domain}}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/{{.Domain}}/privkey.pem;

    location / {
        proxy_pass http://{{.TargetHost}}:{{.TargetPort}};

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
`

func GenerateHTTPProxyConfig(config ProxyConfig) ([]byte, error) {
	return executeTemplate(
		"http-reverse-proxy",
		httpProxyTemplate,
		config,
	)
}

func GenerateHTTPSProxyConfig(config ProxyConfig) ([]byte, error) {
	return executeTemplate(
		"https-reverse-proxy",
		httpsProxyTemplate,
		config,
	)
}

func executeTemplate(
	name string,
	templateContent string,
	config ProxyConfig,
) ([]byte, error) {

	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf(
			"parse nginx template: %w",
			err,
		)
	}

	var content bytes.Buffer

	if err := tmpl.Execute(&content, config); err != nil {
		return nil, fmt.Errorf(
			"generate nginx config: %w",
			err,
		)
	}

	return content.Bytes(), nil
}

func ConfigPath(configDir string, domain string) string {
	filename := "nginx-manager--" + domain + ".conf"

	return filepath.Join(
		configDir,
		filename,
	)
}

func WriteFileAtomic(path string, content []byte) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf(
			"create config directory: %w",
			err,
		)
	}

	tempFile, err := os.CreateTemp(
		dir,
		".nginx-manager-*.tmp",
	)

	if err != nil {
		return fmt.Errorf(
			"create temporary config: %w",
			err,
		)
	}

	tempPath := tempFile.Name()

	defer os.Remove(tempPath)

	if _, err := tempFile.Write(content); err != nil {
		tempFile.Close()

		return fmt.Errorf(
			"write temporary config: %w",
			err,
		)
	}

	if err := tempFile.Chmod(0644); err != nil {
		tempFile.Close()

		return fmt.Errorf(
			"chmod config: %w",
			err,
		)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf(
			"close temporary config: %w",
			err,
		)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf(
			"activate config: %w",
			err,
		)
	}

	return nil
}