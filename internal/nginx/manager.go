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
		Simpan kondisi config sebelumnya
		untuk rollback.
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
		Tulis config baru.
	*/

	if err := WriteFileAtomic(
		path,
		content,
	); err != nil {
		return "", err
	}

	/*
		Test config sebelum reload.
	*/

	if err := TestConfig(); err != nil {

		rollbackErr := rollbackConfig(
			path,
			hadOldConfig,
			oldContent,
		)

		if rollbackErr != nil {
			return "", fmt.Errorf(
				"%v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		return "", err
	}

	/*
		Config valid.
		Sekarang reload nginx.
	*/

	if err := Reload(); err != nil {

		/*
			Kalau reload gagal,
			kembalikan file sebelumnya.
		*/

		rollbackErr := rollbackConfig(
			path,
			hadOldConfig,
			oldContent,
		)

		if rollbackErr != nil {
			return "", fmt.Errorf(
				"%v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		/*
			Pastikan config hasil rollback
			masih valid.
		*/

		if testErr := TestConfig(); testErr != nil {
			return "", fmt.Errorf(
				"%v; rollback config is invalid: %w",
				err,
				testErr,
			)
		}

		/*
			Coba reload ulang menggunakan
			config yang sudah dikembalikan.

			Best effort.
		*/

		if restoreReloadErr := Reload(); restoreReloadErr != nil {
			return "", fmt.Errorf(
				"%v; config rolled back but nginx restore reload failed: %w",
				err,
				restoreReloadErr,
			)
		}

		return "", err
	}

	return path, nil
}

func rollbackConfig(
	path string,
	hadOldConfig bool,
	oldContent []byte,
) error {

	if hadOldConfig {
		if err := WriteFileAtomic(
			path,
			oldContent,
		); err != nil {

			return fmt.Errorf(
				"restore old config: %w",
				err,
			)
		}

		return nil
	}

	/*
		Kalau config sebelumnya tidak ada,
		berarti ini CREATE baru.

		Rollback = hapus.
	*/

	if err := os.Remove(path); err != nil &&
		!os.IsNotExist(err) {

		return fmt.Errorf(
			"remove new config: %w",
			err,
		)
	}

	return nil
}