package docker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	unixSocketPath = "/var/run/docker.sock"   // default UNIX socket location
	npipePath      = `\\.\pipe\docker_engine` // default named pipe on Windows
)

func IsDockerRunning() error {
	httpClient := &http.Client{
		// Set a timeout for the request to the Docker daemon
		Timeout: time.Second * 5,
	}

	// For UNIX-like operating systems (including macOS), use the UNIX socket for Docker
	httpClient.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", unixSocketPath)
		},
	}

	// Perform a GET request to the "/_ping" endpoint of the Docker daemon
	resp, err := httpClient.Get("http://localhost/_ping")
	if err != nil {
		return fmt.Errorf("error pinging the Docker daemon: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Docker daemon is not running; received status code: %d", resp.StatusCode)
	}

	// The status code was HTTP 200 OK; the Docker daemon is running
	return nil
}
