package fluent

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParser(t *testing.T) {
	p := &Parser{}

	logFiles := []string{
		"./test/monitor-1.out.txt",
	}

	for _, file := range logFiles {
		fd, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		data, err := p.Parse(file, fd)
		if err != nil {
			t.Fatal(err)
		}
		if data == nil {
			t.Fatal("data is nil")
		}
		t.Log(spew.Sdump(data))
	}
}
