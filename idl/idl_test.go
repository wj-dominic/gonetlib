package idl

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestIDL(t *testing.T) {
	reqEcho := &REQ_ECHO{Message: "test"}

	fmt.Printf("reqEcho: %v\n", reqEcho)

	out, err := proto.Marshal(reqEcho)
	if err != nil {
		t.Failed()
	}

	fmt.Printf("out: %v\n", out)

	outReqEcho := &REQ_ECHO{}
	if proto.Unmarshal(out, outReqEcho) != nil {
		t.Failed()
		return
	}

	fmt.Printf("outReqEcho: %v\n", outReqEcho)
}
