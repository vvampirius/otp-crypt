package coder

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func Enc(t *testing.T, filePath, otpAuthUrl, text string, key []byte) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	encoder, err := NewEncoder(f, otpAuthUrl, nil, key)
	if err != nil {
		t.Fatal(err.Error())
	}
	if _, err := fmt.Fprintf(encoder, text); err != nil {
		t.Fatal(err.Error())
	}
}

func Dec(t *testing.T, filePath, otpAuthUrl, text string, key []byte) {
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	askKey := func() ([]byte, error) {
		otpKey, _ := otp.NewKeyFromURL(otpAuthUrl)
		passcode, err := totp.GenerateCode(otpKey.Secret(), time.Now())
		if err != nil {
			t.Fatal(err.Error())
		}
		return []byte(passcode), nil
	}

	decoder, err := NewDecoder(f, key, askKey)
	if err != nil {
		t.Fatal(err.Error())
	}

	p, err := ioutil.ReadAll(decoder)
	if err != nil {
		t.Fatal(err.Error())
	}

	s := string(p)
	if s != text {
		ErrorLog.Println(text)
		ErrorLog.Println(s)
		t.Fatal(`don't match!`)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	sKey := `1234`
	otpAuthUrl := `otpauth://totp/ACME%20Co:john.doe@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30`
	text := `Hello World!`

	tmpDir, err := ioutil.TempDir(``, ``)
	if err != nil {
		t.Fatal(err.Error())
	}
	//defer os.RemoveAll(tmpDir)

	filePath := path.Join(tmpDir, `TestEncryptDecrypt`)
	DebugLog.Println(filePath)

	key, err := StringToKey(sKey)
	if err != nil {
		t.Fatal(err.Error())
	}

	Enc(t, filePath, otpAuthUrl, text, key)
	Dec(t, filePath, otpAuthUrl, text, key)

}

func TestWriteOtpImage(t *testing.T) {
	otpAuthUrl := `otpauth://totp/ACME%20Co:john.doe@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30`

	tmpDir, err := ioutil.TempDir(``, ``)
	if err != nil {
		t.Fatal(err.Error())
	}

	filePath := path.Join(tmpDir, `barcode`)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	otpKey, err := otp.NewKeyFromURL(otpAuthUrl)
	if err != nil {
		t.Fatal(err.Error())
	}

	img, err := otpKey.Image(200, 200)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err := png.Encode(f, img); err != nil {
		t.Fatal(err.Error())
	}

	DebugLog.Println(filePath)
}
