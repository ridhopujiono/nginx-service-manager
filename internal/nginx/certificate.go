package nginx

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func CertificateExists(domain string) bool {
	liveDir := os.Getenv("LETSENCRYPT_LIVE_DIR")

	if liveDir == "" {
		liveDir = "/etc/letsencrypt/live"
	}

	certPath := filepath.Join(
		liveDir,
		domain,
		"fullchain.pem",
	)

	keyPath := filepath.Join(
		liveDir,
		domain,
		"privkey.pem",
	)

	if _, err := os.Stat(certPath); err != nil {
		return false
	}

	if _, err := os.Stat(keyPath); err != nil {
		return false
	}

	return true
}

func IssueCertificate(domain string) error {
	mode := os.Getenv("CERTIFICATE_MODE")

	if mode == "" {
		mode = "disabled"
	}

	if mode == "disabled" {
		return fmt.Errorf(
			"certificate issuance is disabled",
		)
	}

	if mode != "certbot" {
		return fmt.Errorf(
			"unsupported certificate mode: %s",
			mode,
		)
	}

	email := strings.TrimSpace(
		os.Getenv("LETSENCRYPT_EMAIL"),
	)

	if email == "" {
		return fmt.Errorf(
			"LETSENCRYPT_EMAIL is not configured",
		)
	}

	webroot := os.Getenv("ACME_WEBROOT")

	if webroot == "" {
		webroot = "/var/lib/nginx-manager/acme"
	}

	if err := os.MkdirAll(
		filepath.Join(
			webroot,
			".well-known",
			"acme-challenge",
		),
		0755,
	); err != nil {

		return fmt.Errorf(
			"create ACME webroot: %w",
			err,
		)
	}

	binary := os.Getenv("CERTBOT_BINARY")

	if binary == "" {
		binary = "/usr/bin/certbot"
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Minute,
	)

	defer cancel()

	args := []string{
		"certonly",

		"--webroot",
		"--webroot-path",
		webroot,

		"--domain",
		domain,

		"--cert-name",
		domain,

		"--non-interactive",
		"--agree-tos",

		"--email",
		email,

		"--keep-until-expiring",
	}

	cmd := exec.CommandContext(
		ctx,
		binary,
		args...,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf(
			"certbot certificate request timed out",
		)
	}

	if err != nil {
		output := strings.TrimSpace(
			stderr.String(),
		)

		if output == "" {
			output = strings.TrimSpace(
				stdout.String(),
			)
		}

		return fmt.Errorf(
			"certbot failed: %s",
			output,
		)
	}

	return nil
}