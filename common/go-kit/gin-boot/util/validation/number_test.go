package validation

import (
	"log"
	"testing"
)

func TestID(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		err := ID(-1)
		log.Println(err)
	})

	t.Run("multiple", func(t *testing.T) {
		err := ID([]int{1, -2}, []int{3, 4})
		log.Println(err)
	})
}
