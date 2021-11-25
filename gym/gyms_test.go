package gym

import (
	"fmt"
	"testing"
)

func TestCreateGyms(t *testing.T){
	gymManager := GetGyms()

	fmt.Println(gymManager)
}