package main

import (
	"fmt"
	"message"
)

func main(){
	msg := message.NewMessage(true)

	msg.Push("test")

	fmt.Println("origin : ", *msg.GetPayloadBuffer())

	msg.Encode()

	fmt.Println("encoded : ", *msg.GetPayloadBuffer())

	msg.Decode()

	fmt.Println("decoded : ", *msg.GetPayloadBuffer())

	var testMsg string

	msg.Pop(&testMsg)

	fmt.Println(testMsg)
}