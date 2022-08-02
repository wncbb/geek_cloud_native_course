package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/kr/pretty"
)

const (
	EnvVersionKey        = "VERSION"
	HttpHeaderVersionKey = "VERSION"
)

func index(resp http.ResponseWriter, req *http.Request) {
	os.Setenv(EnvVersionKey, "v0.0.1")
	resp.Header().Set(HttpHeaderVersionKey, os.Getenv(EnvVersionKey))

	for k, v := range req.Header {
		for _, vv := range v {
			resp.Header().Add(k, vv)
		}
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"msg":"success", "code":0}`))

	fmt.Printf("Request from %s, Resp status %d\n", getClientIP(req), http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/localhost/healthz", index)
	mux.HandleFunc("/healthz", index)
	mux.HandleFunc("/", index)

	if err := http.ListenAndServe(":7878", mux); err != nil {
		panic(err)
	}
	pretty.Println("start")
}

func getClientIP(r *http.Request) string {
	xForwardedForIP := strings.Split(strings.TrimSpace(r.Header.Get("X-Forwarded-For")), ",")[0]
	if xForwardedForIP != "" {
		return xForwardedForIP
	}

	xRealIP := strings.Split(strings.TrimSpace(r.Header.Get("X-Real-Ip")), ",")[0]
	if xRealIP != "" {
		return xRealIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// log err
		return ""
	}

	return host
}
