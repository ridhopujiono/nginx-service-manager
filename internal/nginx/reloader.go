package nginx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Reload() error {
	mode := os.Getenv("NGINX_RELOAD_MODE")

	if mode == "" {
		mode = "noop"
	}

	switch mode {
	case "noop":
		return nil

	case "systemctl":
		return reloadWithSystemctl()

	default:
		return fmt.Errorf(
			"unsupported nginx reload mode: %s",
			mode,
		)
	}
}

func reloadWithSystemctl() error {
	systemctl := os.Getenv("SYSTEMCTL_BINARY")

	if systemctl == "" {
		systemctl = "/usr/bin/systemctl"
	}

	cmd := exec.Command(
		systemctl,
		"reload",
		"nginx",
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())

		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		return fmt.Errorf(
			"nginx reload failed: %s",
			output,
		)
	}

	return nil
}