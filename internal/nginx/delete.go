package nginx

import (
	"errors"
	"fmt"
	"os"
)

var ErrProxyNotFound = errors.New(
	"proxy not found",
)

func DeleteProxy(
	configDir string,
	domain string,
) error {

	if configDir == "" {
		return fmt.Errorf(
			"nginx config directory is not configured",
		)
	}

	path := ConfigPath(
		configDir,
		domain,
	)

	/*
		Backup file sebelum delete.
	*/

	oldContent, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		return ErrProxyNotFound
	}

	if err != nil {
		return fmt.Errorf(
			"read existing config: %w",
			err,
		)
	}

	/*
		Delete config.
	*/

	if err := os.Remove(path); err != nil {
		return fmt.Errorf(
			"delete config: %w",
			err,
		)
	}

	/*
		Pastikan config Nginx masih valid.
	*/

	if err := TestConfig(); err != nil {

		if rollbackErr := WriteFileAtomic(
			path,
			oldContent,
		); rollbackErr != nil {

			return fmt.Errorf(
				"%v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		return err
	}

	/*
		Reload Nginx.
	*/

	if err := Reload(); err != nil {

		/*
			Kembalikan config yang dihapus.
		*/

		if rollbackErr := WriteFileAtomic(
			path,
			oldContent,
		); rollbackErr != nil {

			return fmt.Errorf(
				"%v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		/*
			Pastikan config hasil restore valid.
		*/

		if testErr := TestConfig(); testErr != nil {
			return fmt.Errorf(
				"%v; rollback config invalid: %w",
				err,
				testErr,
			)
		}

		/*
			Restore runtime Nginx juga.
		*/

		if reloadErr := Reload(); reloadErr != nil {
			return fmt.Errorf(
				"%v; config restored but nginx reload failed: %w",
				err,
				reloadErr,
			)
		}

		return err
	}

	return nil
}