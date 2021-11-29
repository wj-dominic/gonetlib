package message

import (
	"fmt"
	"testing"
)

/*
1. 메시지 객체 생성
2. 메시지에 데이터 삽입
3. 메시지에 헤더 세팅
4. 메시지 암호화
5. 메시지 전송
*/

func TestMessage_Push(t *testing.T){
	msg := NewMessage(true)
	
	netHeader := NetHeader{
		PacketType:    1,
		CryptoType:    2,
		RandKey:       3,
		PayloadLength: 4,
		CheckSum:      5,
	}

	fmt.Println("Header : ", netHeader)

	msg.Push(netHeader)

	msg.Push("testslfjlskdjfklsdfhl")

	msg.Push(int16(10))

	fmt.Println("buffer : ", msg.GetBuffer())

	var peekNetHeader NetHeader

	msg.Pop(&peekNetHeader)

	fmt.Println("Header : ", peekNetHeader)

	var test string

	msg.Pop(&test)

	fmt.Println("string : " , test)

	var testInt int
	msg.Pop(&testInt)

	fmt.Println("int : ", testInt)
}

/*
func TestMsgEncodingXOR(t *testing.T) {
	msg := NewMessage(true)

	msg.Push("it is test code")
	msg.Push("next code")

	msg.SetHeader(SYN, XOR)

	fmt.Println("buffer : ", msg.GetBuffer())

	msg.EncodeXOR(0xa9)

	fmt.Println("encoded from xor : ", msg.GetBuffer())

	msg.DecodeXOR(0xa9)

	fmt.Println("decoded from xor : ", msg.GetBuffer())

	var testMsg string
	msg.Pop(&testMsg)
	fmt.Println(testMsg)

	msg.Pop(&testMsg)
	fmt.Println(testMsg)

	if msg == nil{
		t.Fail()
	}
}

func TestMsgEncodingRSA(t *testing.T) {
	msg := NewMessage(true)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fail()
	}

	publicKey := privateKey.PublicKey

	msg.Push("it is test code")
	msg.Push("next code")

	msg.SetHeader(SYN, RSA)

	fmt.Println("buffer : ", msg.GetBuffer())

	msg.EncodeRSA(&publicKey)

	fmt.Println("encoded from rsa : ", msg.GetBuffer())

	msg.DecodeRSA(privateKey)

	fmt.Println("decoded from rsa : ", msg.GetBuffer())

	var testMsg string
	msg.Pop(&testMsg)
	fmt.Println(testMsg)

	msg.Pop(&testMsg)
	fmt.Println(testMsg)

	if msg == nil{
		t.Fail()
	}

}

//*/