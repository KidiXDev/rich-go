//go:build windows
// +build windows

package ipc

import (
	"fmt"
	npipe "gopkg.in/natefinch/npipe.v2"
	"time"
)

// OpenSocket opens the discord-ipc-0 named pipe
func OpenSocket() error {
	var lastErr error
	for i := 0; i < 10; i++ {
		paths := []string{
			fmt.Sprintf(`\\?\pipe\discord-ipc-%d`, i),
			fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, i),
		}
		for _, path := range paths {
			// DialTimeout avoids hanging forever when Discord is not running.
			sock, err := npipe.DialTimeout(path, time.Second*2)
			if err != nil {
				lastErr = err
				continue
			}

			socket = sock
			return nil
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return fmt.Errorf("no discord ipc pipe found")
}
