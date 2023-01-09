package generator

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	generator := NewGenerator(NewIDLParser())
	if generator.Generate("./Idl") == false {
		t.Failed()
	}
}
