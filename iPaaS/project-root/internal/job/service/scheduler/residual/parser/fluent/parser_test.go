package fluent

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	p := &Parser{}

	logFiles := []string{
		"./test/1.txt",
		"./test/2.txt",
		"./test/3.txt",
		"./test/4.txt",
		"./test/5.txt",
		"./test/current_equals_remaining.txt",
		"./test/discontinuity.txt",
	}

	for _, file := range logFiles {
		fd, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		data, err := p.Parse(fd)
		if err != nil {
			t.Fatal(err)
		}
		if data == nil {
			t.Fatal("data is nil")
		}
	}
}
