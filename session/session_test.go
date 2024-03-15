package session_test

import (
	"bytes"
	"testing"
)

func TestSession(t *testing.T) {
	var buffer bytes.Buffer

	buffer.WriteString("test")
	buffer.WriteRune(1234)

	buffer.Truncate(3)

}
