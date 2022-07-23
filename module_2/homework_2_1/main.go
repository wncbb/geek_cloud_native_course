package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
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
}

func getClientIP(r *http.Request) string {
	xForwardedForIPs := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	if len(xForwardedForIPs) >= 1 && xForwardedForIPs[0] != "" {
		return xForwardedForIPs[0]
	}

	xRealIPs := strings.Split(r.Header.Get("X-Real-Ip"), ",")
	if len(xRealIPs) >= 1 && xRealIPs[0] != "" {
		return xRealIPs[0]
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// log err
		return ""
	}

	return host
}
