package ipc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

var socket net.Conn

// Choose the right directory to the ipc socket and return it
func GetIpcPath() string {
	variablesnames := []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

	if _, err := os.Stat("/run/user/1000/snap.discord"); err == nil {
		return "/run/user/1000/snap.discord"
	}

	if _, err := os.Stat("/run/user/1000/.flatpak/com.discordapp.Discord/xdg-run"); err == nil {
		return "/run/user/1000/.flatpak/com.discordapp.Discord/xdg-run"
	}

	for _, variablename := range variablesnames {
		path, exists := os.LookupEnv(variablename)

		if exists {
			return path
		}
	}

	return "/tmp"
}

func CloseSocket() error {
	if socket != nil {
		socket.Close()
		socket = nil
	}
	return nil
}

// Read the socket response
func Read() ([]byte, error) {
	if socket == nil {
		return nil, errors.New("discord ipc socket is not open")
	}

	header := make([]byte, 8)
	if _, err := io.ReadFull(socket, header); err != nil {
		return nil, err
	}

	payloadLength := int(binary.LittleEndian.Uint32(header[4:8]))
	if payloadLength < 0 {
		return nil, fmt.Errorf("invalid payload length: %d", payloadLength)
	}
	if payloadLength == 0 {
		return []byte{}, nil
	}

	payload := make([]byte, payloadLength)
	if _, err := io.ReadFull(socket, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// Send opcode and payload to the unix socket
func Send(opcode int, payload string) ([]byte, error) {
	if socket == nil {
		return nil, errors.New("discord ipc socket is not open")
	}

	buf := make([]byte, 8+len(payload))
	binary.LittleEndian.PutUint32(buf[0:4], uint32(opcode))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(payload)))
	copy(buf[8:], payload)

	if _, err := socket.Write(buf); err != nil {
		return nil, err
	}

	return Read()
}
