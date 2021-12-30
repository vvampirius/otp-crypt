package main

import "errors"

var ErrKeyTooLong = errors.New(`Key length more than 32 bytes`)

func StringToKey(s string) ([]byte, error) {
	sBytes := []byte(s)
	if len(sBytes) > 32 {
		ErrorLog.Println(ErrKeyTooLong.Error())
		return nil, ErrKeyTooLong
	}
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = byte(i + 1)
	}
	copy(key, sBytes)
	return key, nil
}
