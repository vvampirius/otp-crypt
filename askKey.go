package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var ErrEmptyInput = errors.New(`Empty input`)

type AskKey struct {
	EnvName string
	Cmd string
	AllowStdin bool
	StdinPrompt string
}

func (askKey *AskKey) GetKey() ([]byte, error) {
	if askKey.EnvName != `` {
		if envContent := os.Getenv(askKey.EnvName); envContent != `` {
			return []byte(envContent), nil
		}
	}
	if askKey.Cmd != `` {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Getenv(`SHELL`), `-c`, askKey.Cmd)
		output, err := cmd.Output()
		return bytes.TrimSuffix(output, []byte("\n")), err
	}
	if askKey.AllowStdin {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(askKey.StdinPrompt)
		input, err := reader.ReadBytes(byte('\n'))
		if err != nil { ErrorLog.Println(err.Error()) }
		input = bytes.TrimSuffix(input, []byte("\n"))
		return input, err
	}
	return []byte{}, ErrEmptyInput
}