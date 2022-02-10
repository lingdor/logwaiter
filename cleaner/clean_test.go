package cleaner

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFiles(t *testing.T) {
	files, err := filepath.Glob("/root/logwaiter/log/*")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)

	}
}
