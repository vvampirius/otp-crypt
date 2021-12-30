package coder

import (
	"errors"
	"log"
	"os"
)

const formatVersion = 0x00

var (
	DebugLog = log.New(os.Stderr, `debug#`, log.Lshortfile)
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
	ErrBadInput = errors.New(`Bad input`)
)
