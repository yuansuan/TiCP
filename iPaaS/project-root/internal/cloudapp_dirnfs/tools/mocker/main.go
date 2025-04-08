package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
)

func generate(args []string) {
	digits := regexp.MustCompile("^\\d+$")
	for _, arg := range args {
		if digits.MatchString(arg) {
			if n, err := strconv.ParseInt(arg, 10, 64); err == nil {
				if id, err := hashid.Encode(snowflake.ID(n)); err == nil {
					fmt.Printf("%20d :\t%s\n", n, id)
				}
			}
		} else if len(arg) == 34 && strings.HasPrefix(arg, "YS") {
			if id, err := hashid.Decode(arg[2:]); err == nil {
				fmt.Printf("%s :\t%s\n", arg, id)
			}
		} else {
			var sb strings.Builder
			for _, r := range []rune(arg) {
				if unicode.IsDigit(r) || unicode.IsLetter(r) {
					sb.WriteRune(r)
				}
			}

			if n, err := snowflake.ParseString(sb.String()); err == nil {
				if id, err := hashid.Encode(n); err == nil {
					fmt.Printf("%20s :\t%s\n", sb.String(), id)
				}
			}
		}
	}
}

func register(args []string) {
	var (
		server   string
		username int64
		userhome string
	)

	flag.StringVar(&server, "server", "", "server address")
	flag.Int64Var(&username, "username", 1234, "username for user")
	flag.StringVar(&userhome, "userhome", "", "user home dir")
	if err := flag.CommandLine.Parse(args); err != nil {
		panic(err)
	}

	uid, err := hashid.Encode(snowflake.ID(username))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	buf.WriteString(uid)

	body := map[string]string{"sub_path": userhome}
	bs, _ := json.Marshal(body)

	crc := make([]byte, 4)
	binary.BigEndian.PutUint32(crc, crc32.ChecksumIEEE(bs))
	buf.Write(crc)
	buf.Write(bs)

	resp, err := http.Post(server+"/"+uid, "application/json", &buf)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	fmt.Printf("Register: %d - %s\n", resp.StatusCode, resp.Status)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mocker [generate|register] args...")
		return
	}

	switch os.Args[1] {
	case "generate":
		generate(os.Args[2:])
	case "register":
		register(os.Args[2:])
	}
}
