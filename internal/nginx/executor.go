package nginx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func TestConfig() error {
	binary := os.Getenv("NGINX_BINARY")

	if binary == "" {
		binary = "/usr/sbin/nginx"
	}

	testConfig := os.Getenv("NGINX_TEST_CONFIG")

	args := []string{"-t"}

	if testConfig != "" {
		args = append(
			args,
			"-c",
			testConfig,
		)
	}

	cmd := exec.Command(binary, args...)

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
			"nginx config test failed: %s",
			output,
		)
	}

	return nil
}