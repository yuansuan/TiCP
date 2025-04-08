package registry

import (
	"testing"
)

type MyFloat float64

type Book struct {
	name string
	time string
}

func (b *Book) show() (resp string) {
	return "book-show"
}

func Test_Cache(t *testing.T) {
	reg := GetRegistry()
	book := new(Book)
	t.Log(reg.Set("a", 1.1))
	t.Log(reg.Set("b", book))
	t.Log(reg.Set("c", nil))
	t.Log(reg.Get("a"))
	t.Log("b.show: " + reg.Get("b").(*Book).show())
	t.Log(reg.Get("c"))
	t.Log("delete a")
	reg.Delete("a")
	t.Log(reg.Get("a"))
	reg.Clear()
	t.Log("clear .")
	t.Log(reg.Get("b"))

}
