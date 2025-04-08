package xhttp

import "net/http"

func AppendQuery(req *http.Request, key, val string) *http.Request {
	qs := req.URL.Query()
	qs.Set(key, val)
	req.URL.RawQuery = qs.Encode()
	return req
}

func AddHeader(req *http.Request, key, val string) *http.Request {
	req.Header.Add(key, val)
	return req
}
