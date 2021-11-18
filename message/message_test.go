package message_test

import (
	"fmt"
	"message"
	"testing"
)

func TestMessage(t *testing.T) {
	msg := message.NewMessage(true)

	msg.Push("test")

	msg.SetHeader(message.SYN, message.XOR)

	fmt.Println("buffer : ", msg.GetBuffer())

	var testMsg string

	msg.Pop(&testMsg)

	fmt.Println(testMsg)

	if msg == nil{
		t.Fail()
	}
}