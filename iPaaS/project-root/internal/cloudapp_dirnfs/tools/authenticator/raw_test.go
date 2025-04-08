package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestAdduser(t *testing.T) {
	resp, err := http.Post("http://10.0.1.123:8081/05a3a41adaf71e0ec1dd940bec79afed", "application/sharing", strings.NewReader("0103bc6fabfe5f8e48164ec011b67304"))
	if err != nil {
		println(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		println(resp.StatusCode)
	}
}
