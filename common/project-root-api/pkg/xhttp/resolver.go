package xhttp

import (
	"net/http"
)

type ResponseResolver func(resp *http.Response) error
