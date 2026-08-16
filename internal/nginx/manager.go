package nginx

import (
	"fmt"
	"os"
)

func CreateOrUpdateProxy(
	configDir string,
	config ProxyConfig,
) (string, error) {

	if configDir == "" {
		return "", fmt.Errorf(
			"nginx config directory is not configured",
		)
	}

	content, err := GenerateProxyConfig(config)

	if err != nil {
		return "", err
	}

	path := ConfigPath(
		configDir,
		config.Domain,
	)

	/*
		Backup config lama.
	*/

	var oldContent []byte
	hadOldConfig := false

	if existing, err := os.ReadFile(path); err == nil {

		oldContent = existing
		hadOldConfig = true

	} else if !os.IsNotExist(err) {

		return "", fmt.Errorf(
			"read existing config: %w",
			err,
		)
	}

	/*
		Write config baru.
	*/

	if err := WriteFileAtomic(path, content); err != nil {
		return "", err
	}

	/*
		Test seluruh konfigurasi nginx.
	*/

	if err := TestConfig(); err != nil {

		/*
			Rollback.
		*/

		if hadOldConfig {

			rollbackErr := WriteFileAtomic(
				path,
				oldContent,
			)

			if rollbackErr != nil {
				return "", fmt.Errorf(
					"%v; rollback failed: %w",
					err,
					rollbackErr,
				)
			}

		} else {

			if removeErr := os.Remove(path); removeErr != nil {
				return "", fmt.Errorf(
					"%v; rollback failed: %w",
					err,
					removeErr,
				)
			}
		}

		return "", err
	}

	return path, nil
}