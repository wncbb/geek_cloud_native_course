package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kr/pretty"
)

const (
	EnvVersionKey        = "VERSION"
	HttpHeaderVersionKey = "VERSION"
)

type Data struct {
	Company string
}

type Resp struct {
	Code    int
	Version string
	Time    time.Time
	Message string
	Data    Data
}

func index(resp http.ResponseWriter, req *http.Request) {
	version := os.Getenv(EnvVersionKey)
	resp.Header().Set(HttpHeaderVersionKey, version)

	for k, v := range req.Header {
		for _, vv := range v {
			resp.Header().Add(k, vv)
		}
	}

	r := &Resp{
		Time:    time.Now().UTC(),
		Version: version,
	}

	company, err := os.ReadFile("/data/company")
	if err != nil {
		r.Code = 500
		r.Message = err.Error()
		output, _ := json.Marshal(r)
		resp.Write([]byte(output))
		return
	}

	r.Code = 200
	r.Message = "success"
	r.Data = Data{
		Company: string(company),
	}
	output, _ := json.Marshal(r)

	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte(output))

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
