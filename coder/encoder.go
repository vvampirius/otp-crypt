package coder

import (
	"io"

	"github.com/pquerna/otp"
	"github.com/vvampirius/mygolibs/encryption"
	"github.com/vvampirius/mygolibs/rwblocks"
)

type Encoder struct {
	writeFd io.Writer
	cryptor *encryption.CFB
}

func (encoder *Encoder) Write(p []byte) (int, error) {
	if err := rwblocks.Write(encoder.writeFd, encoder.cryptor.Encrypt(p)); err != nil {
		ErrorLog.Println(err.Error())
		return 0, err
	}
	return len(p), nil
}

func (encoder *Encoder) WriteVersion() error {
	if _, err := encoder.writeFd.Write([]byte{formatVersion}); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	return nil
}

func (encoder *Encoder) WriteIV() error {
	if err := rwblocks.Write(encoder.writeFd, encoder.cryptor.IV); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	return nil
}

func (encoder *Encoder) WriteOtpAuthUrl(otpAuthUrl string, key []byte) error {
	otpAuthUrlStruct, err := otp.NewKeyFromURL(otpAuthUrl)
	if err != nil {
		ErrorLog.Println()
		return err
	}

	encrypted := encoder.cryptor.Encrypt([]byte(otpAuthUrl))
	if err := rwblocks.Write(encoder.writeFd, encrypted); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	dataKey := DataKey(key, encoder.cryptor.IV, otpAuthUrlStruct.Secret())
	if err := encoder.cryptor.SetKey(dataKey); err != nil {
		ErrorLog.Println(err.Error())
		return err
	}

	return nil
}

func NewEncoder(w io.Writer, otpAuthUrl string, iv, key []byte) (*Encoder, error) {
	if otpAuthUrl == `` && w == nil {
		ErrorLog.Println(ErrBadInput)
		return nil, ErrBadInput
	}

	cfb, err := encryption.NewCFB(iv, key)
	if err != nil {
		return nil, err
	}

	encoder := Encoder{
		writeFd: w,
		cryptor: cfb,
	}
	if err := encoder.WriteVersion(); err != nil {
		return nil, err
	}
	if err := encoder.WriteIV(); err != nil {
		return nil, err
	}
	if err := encoder.WriteOtpAuthUrl(otpAuthUrl, key); err != nil {
		return nil, err
	}

	return &encoder, nil
}
