package singleton

import (
	"fmt"
	"testing"
)

type Book struct {
	ID     uint32
	Title  string
	Author string
}

func (book *Book) Init() {
	book.ID = 0
	book.Title = "Empty"
	book.Author = "Empty"
}

func TestSingleton_SetInstance(t *testing.T) {
	book := GetInstance[Book]()
	if book == nil {
		t.Failed()
	}

	fmt.Println(book)

	book.ID = 1
	book.Title = "love ballard"
	book.Author = "kim"

	fmt.Println(book)

	for i := 0; i < 10; i++ {
		v := GetInstance[Book]()

		if v != book {
			t.Failed()
		}
	}
}

func TestSingleton_GetInstance(t *testing.T) {
	book := GetInstance[Book]()

	fmt.Println(book)
}
