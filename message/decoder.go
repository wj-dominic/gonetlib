package message

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
)

type XORDecoder struct {
}

func NewXORDecoder() *XORDecoder {
	return &XORDecoder{}
}

func (dc *XORDecoder) Decode(key interface{}, buf []byte) bool {
	numberKey := key.(uint32)

	randKey := buf[0]
	dstBuffer := buf[1:]

	num := uint8(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint8(numberKey)+num)
		dstBuffer[i] = p ^ uint8(randKey+num)
	}

	return true
}

type RSADecoder struct {
}

func NewRSADecoder() *RSADecoder {
	return &RSADecoder{}
}

func (dc *RSADecoder) Decode(key interface{}, buf []byte) bool {
	privateKey := key.(*rsa.PrivateKey)

	plainMsg, err := rsa.DecryptPKCS1v15(cryptoRand.Reader, privateKey, buf)
	if err != nil {
		return false
	}

	buf = plainMsg

	return true
}
