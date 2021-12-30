package main

import (
	"fmt"
	"github.com/vvampirius/otp-crypt/coder"
	"io"
	"os"
)

func Encrypt(srcFilePath, dstFilePath, askKeyCmd string, denyConsoleInteractive bool) error {
	askOtpAuthUrl := AskKey{
		EnvName: `OTP_AUTH_URL`,
	}
	if !denyConsoleInteractive {
		askOtpAuthUrl.AllowStdin = true
		askOtpAuthUrl.StdinPrompt = `Enter OTP auth URL: `
	}
	pOtpAuthUrl, err := askOtpAuthUrl.GetKey()
	if err != nil {
		fmt.Fprintln(os.Stderr, `You can provide OTP auth URL by environment variable OTP_AUTH_URL`)
		return err
	}
	otpAuthUrl := string(pOtpAuthUrl)

	askKey := AskKey{
		EnvName: `ENCRYPTION_KEY`,
		Cmd: askKeyCmd,
	}
	if !denyConsoleInteractive {
		askKey.AllowStdin = true
		askKey.StdinPrompt = `Enter encryption key: `
	}
	key, err := askKey.GetKey()
	if err != nil {
		fmt.Fprintln(os.Stderr, `You can provide encryption key by environment variable ENCRYPTION_KEY`)
		return err
	}

	var src io.Reader
	src = os.Stdin
	if srcFilePath != `` {
		fSrc, err := os.Open(srcFilePath)
		if err != nil {
			ErrorLog.Println(err.Error())
			return err
		}
		defer fSrc.Close()
		src = fSrc
	}

	var dst io.Writer
	dst = os.Stdout
	if dstFilePath != `` {
		fDst, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			ErrorLog.Println(err.Error())
			return err
		}
		defer fDst.Close()
		dst = fDst
	}

	encoder, err := coder.NewEncoder(dst, otpAuthUrl, nil, key)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	if _, err := io.Copy(encoder, src); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	return nil
}
