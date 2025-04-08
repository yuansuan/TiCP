package job

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateImageNameMap(t *testing.T) {
	testCases := []struct {
		imageNameMap map[string][]string
		regName      string
		fileName     string
		expectedMap  map[string][]string
	}{
		{
			imageNameMap: make(map[string][]string),
			regName:      "example",
			fileName:     "image.png",
			expectedMap: map[string][]string{
				"example": {"image.png"},
			},
		},
		{
			imageNameMap: map[string][]string{
				"example": {"image.png"},
			},
			regName:  "example",
			fileName: "image2.png",
			expectedMap: map[string][]string{
				"example": {"image.png", "image2.png"},
			},
		},
		{
			imageNameMap: make(map[string][]string),
			regName:      "example2",
			fileName:     "image.png",
			expectedMap: map[string][]string{
				"example2": {"image.png"},
			},
		},
		{
			imageNameMap: map[string][]string{
				"example2": {"image.png"},
			},
			regName:  "example2",
			fileName: "image.png",
			expectedMap: map[string][]string{
				"example2": {"image.png", "image.png"},
			},
		},
	}

	for _, tc := range testCases {
		UpdateImageNameMap(tc.imageNameMap, tc.regName, tc.fileName)
		assert.Equal(t, tc.expectedMap, tc.imageNameMap)
	}
}

type JobType string

const (
	JobTypeA JobType = "A"
	JobTypeB JobType = "B"
)

type A struct{}

func (a *A) Func(name string) string {
	return "A" + name
}

func newA() *A {
	return &A{}
}

type B struct{}

func (b *B) Func(name string) string {
	return name + "B"
}

func newB() *B {
	return &B{}
}

func SomeFunc(names []string) {
	var fn func(string) string

	jobType := getJob()

	if jobType == JobTypeA {
		fn = newA().Func
	} else {
		fn = newB().Func
	}

	AnotherFunc(fn, names)
}

func getJob() (jobType JobType) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	if r.Intn(100)%2 == 1 {
		jobType = JobTypeA
	} else {
		jobType = JobTypeB
	}
	return
}

func AnotherFunc(fn func(string) string, names []string) string {
	for _, name := range names {
		PrintFunc(fn(name), "hello ")
	}
	return ""
}

func PrintFunc(name string, action string) {
	fmt.Println(action + name)
}

func TestSome(t *testing.T) {
	names := []string{"jack", "tom", "jerry"}

	SomeFunc(names)
}
