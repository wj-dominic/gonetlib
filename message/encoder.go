package message

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
)

type XOREncoder struct {
}

func NewXOREncoder() *XOREncoder {
	return &XOREncoder{}
}

func (ec *XOREncoder) Encode(key interface{}, buf []byte) bool {
	numberKey := key.(uint32)

	randKey := buf[0]
	dstBuffer := buf[1:]

	num := uint8(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(randKey+num)
		dstBuffer[i] = p ^ uint8(uint8(numberKey)+num)
	}

	return true
}

type RSAEncoder struct {
}

func (ec *RSAEncoder) Encode(key interface{}, buf []byte) bool {
	publicKey := key.(*rsa.PublicKey)

	cipherMsg, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, publicKey, buf)
	if err != nil {
		return false
	}

	buf = cipherMsg

	return true
}
