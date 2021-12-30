package coder

import (
	"errors"
	"io"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/vvampirius/mygolibs/encryption"
	"github.com/vvampirius/mygolibs/rwblocks"
)

var ErrOTPpassCodeNotValid = errors.New(`OTP passcode is not valid`)

type Decoder struct {
	readFd      io.Reader
	cryptor     *encryption.CFB
	otpAskFunc  func() ([]byte, error)
	otpKey      *otp.Key
	unreadBlock []byte
}

func (decoder *Decoder) readOtpKey() error {
	encrypted, err := rwblocks.Read(decoder.readFd)
	if err != nil {
		return err
	}
	otpUrlBytes := decoder.cryptor.Decrypt(encrypted)
	otpKey, err := otp.NewKeyFromURL(string(otpUrlBytes))
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	decoder.otpKey = otpKey
	return nil
}

func (decoder *Decoder) ValidateOTP() error {
	passcode, err := decoder.otpAskFunc()
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	if !totp.Validate(string(passcode), decoder.otpKey.Secret()) {
		ErrorLog.Println(ErrOTPpassCodeNotValid)
		return ErrOTPpassCodeNotValid
	}
	return nil
}

// ReadBlock returns decrypted block from decoder io.Reader
func (decoder *Decoder) ReadBlock() ([]byte, error) {
	encrypted, err := rwblocks.Read(decoder.readFd)
	if err != nil {
		return nil, err
	}
	return decoder.cryptor.Decrypt(encrypted), nil
}

func (decoder *Decoder) Read(p []byte) (int, error) {
	dstLength := len(p)
	var i int
	for {
		if i == dstLength {
			break
		}
		if decoder.unreadBlock == nil || len(decoder.unreadBlock) == 0 {
			block, err := decoder.ReadBlock()
			if block == nil || len(block) == 0 {
				decoder.unreadBlock = nil
				return i, err
			}
			decoder.unreadBlock = block
		}
		p[i] = decoder.unreadBlock[0]
		i++
		decoder.unreadBlock = decoder.unreadBlock[1:]
	}

	return i, nil
}

func NewDecoder(r io.Reader, key []byte, otpAskFunc func() ([]byte, error)) (*Decoder, error) {
	if r == nil && otpAskFunc == nil {
		ErrorLog.Println(ErrBadInput)
		return nil, ErrBadInput
	}

	// Read first byte with version. It can be used in the future.
	if _, err := r.Read([]byte{0}); err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}

	iv, err := rwblocks.Read(r)
	if err != nil {
		return nil, err
	}

	cfb, err := encryption.NewCFB(iv, key)
	if err != nil {
		return nil, err
	}

	decoder := Decoder{
		readFd:     r,
		cryptor:    cfb,
		otpAskFunc: otpAskFunc,
	}

	if err := decoder.readOtpKey(); err != nil {
		return nil, err
	}
	if err := decoder.ValidateOTP(); err != nil {
		return nil, err
	}

	dataKey := DataKey(key, decoder.cryptor.IV, decoder.otpKey.Secret())
	if err := decoder.cryptor.SetKey(dataKey); err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}

	return &decoder, nil
}
