package logwriter

import (
	"fmt"
	"testing"
)

func TestDate(t *testing.T) {

	month := 2

	fmt.Printf("%2d\n", month)
	fmt.Printf("%02d\n", month)
}
