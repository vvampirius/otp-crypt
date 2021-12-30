package main

import (
	"github.com/vvampirius/otp-crypt/coder"
	"io"
	"os"
)

func Decrypt(srcFilePath, dstFilePath, askKeyCmd, askOtpCmd string, denyConsoleInteractive bool) error {
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
		ErrorLog.Println(err.Error())
		return err
	}

	askOtp := AskKey{
		EnvName: ``,
		Cmd: askOtpCmd,
	}
	if !denyConsoleInteractive {
		askOtp.AllowStdin = true
		askOtp.StdinPrompt = `Enter OTP key: `
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

	decoder, err := coder.NewDecoder(src, key, askOtp.GetKey)
	if err != nil { return err }

	if _, err := io.Copy(dst, decoder); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	return nil
}
