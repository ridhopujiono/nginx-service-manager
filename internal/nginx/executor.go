package nginx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func TestConfig() error {
	binary := os.Getenv("NGINX_BINARY")

	if binary == "" {
		binary = "/usr/sbin/nginx"
	}

	testConfig := os.Getenv("NGINX_TEST_CONFIG")

	var cmd *exec.Cmd

	if testConfig != "" {
		cmd = exec.Command(
			binary,
			"-t",
			"-c",
			testConfig,
		)
	} else {
		cmd = exec.Command(
			binary,
			"-t",
		)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"nginx config test failed: %s",
			stderr.String(),
		)
	}

	return nil
}