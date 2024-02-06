package tests

import (
	"fmt"
	"testing"
)

func TestHello(t *testing.T) {
	fmt.Println("hello")
	if false {
		t.Errorf("err")
	}
}
