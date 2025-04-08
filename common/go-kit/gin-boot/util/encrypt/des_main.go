package main

import (
	"fmt"
	"os"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

// Encrypt Tool
// It will be built into an executable file named encrypt.
func main() {
	args := os.Args
	if len(args) == 2 {
		plainText := os.Args[1]
		cipherText, err := util.Encrypt(plainText)
		if err != nil {
			fmt.Println("")
		} else {
			fmt.Println(cipherText)
		}
	} else {
		fmt.Printf("Invalid syntax\n"+
			"Usage: %v plaintext\n", args[0])
	}
}
