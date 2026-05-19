//go:build !windows
// +build !windows

package ipc

import (
	"fmt"
	"net"
	"time"
)

// OpenSocket opens the discord-ipc-0 unix socket
func OpenSocket() error {
	basePath := GetIpcPath()
	var lastErr error

	for i := 0; i < 10; i++ {
		sock, err := net.DialTimeout("unix", fmt.Sprintf("%s/discord-ipc-%d", basePath, i), time.Second*2)
		if err != nil {
			lastErr = err
			continue
		}

		socket = sock
		return nil
	}

	if lastErr != nil {
		return lastErr
	}

	return fmt.Errorf("no discord ipc socket found in %s", basePath)
}
