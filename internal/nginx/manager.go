package nginx

import (
	"fmt"
	"os"
	"sync"
)

var mutationMu sync.Mutex

func CreateOrUpdateProxy(
	configDir string,
	config ProxyConfig,
) (string, error) {

	mutationMu.Lock()
	defer mutationMu.Unlock()

	if configDir == "" {
		return "", fmt.Errorf(
			"nginx config directory is not configured",
		)
	}

	if config.ACMEWebroot == "" {
		config.ACMEWebroot = os.Getenv(
			"ACME_WEBROOT",
		)

		if config.ACMEWebroot == "" {
			config.ACMEWebroot =
				"/var/lib/nginx-manager/acme"
		}
	}

	path := ConfigPath(
		configDir,
		config.Domain,
	)

	/*
		Simpan kondisi sebelumnya.
	*/

	var oldContent []byte
	hadOldConfig := false

	existing, err := os.ReadFile(path)

	if err == nil {
		oldContent = existing
		hadOldConfig = true
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf(
			"read existing config: %w",
			err,
		)
	}

	/*
		Tidak menggunakan SSL.
		Cukup generate HTTP.
	*/

	if !config.SSL {
		content, err :=
			GenerateHTTPProxyConfig(config)

		if err != nil {
			return "", err
		}

		if err := applyConfig(
			path,
			content,
		); err != nil {

			return "", rollbackAfterFailure(
				err,
				path,
				hadOldConfig,
				oldContent,
			)
		}

		return path, nil
	}

	/*
		SSL requested.

		Kalau certificate sudah ada,
		tidak perlu request Let's Encrypt lagi.
	*/

	if CertificateExists(config.Domain) {
		content, err :=
			GenerateHTTPSProxyConfig(config)

		if err != nil {
			return "", err
		}

		if err := applyConfig(
			path,
			content,
		); err != nil {

			return "", rollbackAfterFailure(
				err,
				path,
				hadOldConfig,
				oldContent,
			)
		}

		return path, nil
	}

	/*
		Certificate belum ada.

		Tahap 1:
		install HTTP bootstrap config
		supaya ACME challenge bisa diakses.
	*/

	httpContent, err :=
		GenerateHTTPProxyConfig(config)

	if err != nil {
		return "", err
	}

	if err := applyConfig(
		path,
		httpContent,
	); err != nil {

		return "", rollbackAfterFailure(
			err,
			path,
			hadOldConfig,
			oldContent,
		)
	}

	/*
		Tahap 2:
		Request Let's Encrypt certificate.
	*/

	if err := IssueCertificate(
		config.Domain,
	); err != nil {

		return "", rollbackAfterFailure(
			err,
			path,
			hadOldConfig,
			oldContent,
		)
	}

	/*
		Tahap 3:
		Certificate sudah tersedia.

		Generate HTTPS config.
	*/

	httpsContent, err :=
		GenerateHTTPSProxyConfig(config)

	if err != nil {
		return "", rollbackAfterFailure(
			err,
			path,
			hadOldConfig,
			oldContent,
		)
	}

	if err := applyConfig(
		path,
		httpsContent,
	); err != nil {

		return "", rollbackAfterFailure(
			err,
			path,
			hadOldConfig,
			oldContent,
		)
	}

	return path, nil
}

func applyConfig(
	path string,
	content []byte,
) error {

	if err := WriteFileAtomic(
		path,
		content,
	); err != nil {

		return err
	}

	if err := TestConfig(); err != nil {
		return err
	}

	if err := Reload(); err != nil {
		return err
	}

	return nil
}

func rollbackAfterFailure(
	originalErr error,
	path string,
	hadOldConfig bool,
	oldContent []byte,
) error {

	if err := rollbackConfig(
		path,
		hadOldConfig,
		oldContent,
	); err != nil {

		return fmt.Errorf(
			"%v; rollback failed: %w",
			originalErr,
			err,
		)
	}

	if err := TestConfig(); err != nil {
		return fmt.Errorf(
			"%v; rollback config invalid: %w",
			originalErr,
			err,
		)
	}

	if err := Reload(); err != nil {
		return fmt.Errorf(
			"%v; config rolled back but reload failed: %w",
			originalErr,
			err,
		)
	}

	return originalErr
}

func rollbackConfig(
	path string,
	hadOldConfig bool,
	oldContent []byte,
) error {

	if hadOldConfig {
		return WriteFileAtomic(
			path,
			oldContent,
		)
	}

	if err := os.Remove(path); err != nil &&
		!os.IsNotExist(err) {

		return fmt.Errorf(
			"remove new config: %w",
			err,
		)
	}

	return nil
}