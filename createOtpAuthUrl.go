package main

import (
	"bufio"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"os"
	"strings"
)

func WriteQrFile(filePath, otpAuthUrl string) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	defer f.Close()

	key, err := otp.NewKeyFromURL(otpAuthUrl)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	img, err := key.Image(200, 200)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	if err := png.Encode(f, img); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	return nil
}

func CreateOtpAuthUrl(qrcodeFilePath string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Issuer: ")
	issuer, _ := reader.ReadString('\n')
	issuer = strings.TrimRight(issuer, "\n")
	fmt.Print("Enter account name: ")
	accountName, _ := reader.ReadString('\n')
	accountName = strings.TrimRight(accountName, "\n")

	opts := totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	}
	key, err := totp.Generate(opts)
	if err !=nil {
		ErrorLog.Println(err.Error())
		return err
	}
	fmt.Println(key)

	if qrcodeFilePath != `` {
		if err := WriteQrFile(qrcodeFilePath, key.String()); err != nil { return err }
	}

	return nil
}
