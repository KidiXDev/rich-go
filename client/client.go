package client

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/KidiXDev/rich-go/ipc"
)

var logged bool

type rpcResponse struct {
	Cmd   string          `json:"cmd"`
	Evt   string          `json:"evt"`
	Nonce string          `json:"nonce"`
	Data  json.RawMessage `json:"data"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Login sends a handshake in the socket and returns an error or nil
func Login(clientid string) error {
	if !logged {
		payload, err := json.Marshal(Handshake{1, clientid})
		if err != nil {
			return err
		}

		err = ipc.OpenSocket()
		if err != nil {
			return err
		}

		response, err := ipc.Send(0, string(payload))
		if err != nil {
			return err
		}
		if err := parseRPCError(response); err != nil {
			_ = ipc.CloseSocket()
			return err
		}
	}
	logged = true

	return nil
}

func Logout() {
	logged = false

	err := ipc.CloseSocket()
	if err != nil {
		panic(err)
	}
}

func SetActivity(activity Activity) error {
	if !logged {
		return nil
	}

	payload, err := json.Marshal(Frame{
		"SET_ACTIVITY",
		Args{
			os.Getpid(),
			mapActivity(&activity),
		},
		getNonce(),
	})

	if err != nil {
		return err
	}

	response, err := ipc.Send(1, string(payload))
	if err != nil {
		return err
	}

	return parseRPCError(response)
}

func getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func parseRPCError(payload []byte) error {
	if len(payload) == 0 {
		return nil
	}

	var response rpcResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		return err
	}
	if response.Evt != "ERROR" {
		return nil
	}

	var rpcErr rpcError
	if err := json.Unmarshal(response.Data, &rpcErr); err != nil {
		return err
	}
	if rpcErr.Message == "" {
		return errors.New("discord rpc returned an unknown error")
	}

	return fmt.Errorf("discord rpc error %d: %s", rpcErr.Code, rpcErr.Message)
}
