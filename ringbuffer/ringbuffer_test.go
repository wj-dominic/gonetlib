package ringbuffer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)


const (
	TESTING_MSG string = "1234567890 abcdefghijklmnopqrstuvwxyz 1234567890 abcdefghijklmnopqrstuvwxyz 12345 abcdefghijklmnopqrstuvwxyz @@@@@@@@@@ "
	buff_size uint32 = 300
)

func TestRingbuffer(t *testing.T) {
	ringBuffer := NewRingBuffer(true, 100)

	message := "it is test code test test test"

	fmt.Println(message)

	writeLength := ringBuffer.Write(message)
	if writeLength != uint32(len(message)) {
		t.Fail()
	}

	var testMsg string
	readLength := ringBuffer.Read(&testMsg, uint32(len(message)))
	if readLength != uint32(len(message)){
		t.Fail()
	}

	fmt.Println(testMsg)
}

func TestInfiniteTest(t *testing.T){
	//무한 링버퍼 삽입 추출 테스트
	ringBuffer := NewRingBuffer(true, 200)

	msgBuffer := []byte(TESTING_MSG)
	msgLength := uint32(len(msgBuffer))

	lastPos := uint32(0)
	pushSize := uint32(0)
	totalPushSize := uint32(0)

	for {
		if msgLength <= totalPushSize {
			totalPushSize = 0
			lastPos = 0
			fmt.Println()
		}

		pushSize = uint32(rand.Intn(int(msgLength - totalPushSize)) + 1)
		totalPushSize += pushSize

		nextPos := lastPos + pushSize
		pushData := msgBuffer[lastPos : nextPos]

		lastPos = nextPos

		ringBuffer.Write(pushData)

		useSize := ringBuffer.GetUseSize()

		var text string
		ringBuffer.Read(&text, useSize)

		fmt.Print(text)

		time.Sleep(5)
	}

}