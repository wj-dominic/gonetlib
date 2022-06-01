package generator

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	generator := NewGenerator()
	generator.Generate("./Idl")
}
