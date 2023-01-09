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

type Test struct {
	Id       int
	Name     string
	Age      uint32
	National string
}

func TestMessage_Push(t *testing.T) {
	msg := NewMessage().LittleEndian().Encoder(NewXOREncoder()).Decoder(NewXORDecoder())

	person := Test{
		Id:       1,
		Name:     "John",
		Age:      30,
		National: "Seoul",
	}

	msg.Push(person)
	msg.Push("testslfjlskdjfklsdfhl")
	msg.Push(10)
	msg.Push(uint(19))
	msg.Push(int8(100))
	msg.Push(byte(90))
	msg.Push(true)
	msg.Push([]byte{1, 2, 3, 4, 5})

	fmt.Println("after push buffer : ", msg.GetPayloadBuffer(), msg.GetPayloadSize())

	var peekStructure Test
	var peekString string
	var peekInt int
	var peekUint uint
	var peekInt8 int8
	var peekByte byte
	var peekBool bool
	var peekBytes []byte

	msg.Pop(&peekStructure)
	msg.Pop(&peekString)
	msg.Pop(&peekInt)
	msg.Pop(&peekUint)
	msg.Pop(&peekInt8)
	msg.Pop(&peekByte)
	msg.Pop(&peekBool)
	msg.Pop(&peekBytes)

	fmt.Println(peekStructure)
	fmt.Println(peekString)
	fmt.Println(peekInt)
	fmt.Println(peekUint)
	fmt.Println(peekInt8)
	fmt.Println(peekByte)
	fmt.Println(peekBool)
	fmt.Println(peekBytes)

	fmt.Println("after pop buffer : ", msg.GetPayloadBuffer(), msg.GetPayloadSize())
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
