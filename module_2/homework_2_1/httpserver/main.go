package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kr/pretty"
	"go.uber.org/zap"
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

func errorHandle(resp http.ResponseWriter, req *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := &Resp{
		Code: http.StatusInternalServerError,
		Time: time.Now().UTC(),
	}

	output, _ := json.Marshal(r)

	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(output))

	logger.Error("Error",
		zap.String("url", req.URL.Path),
		zap.String("msg", "err"),
		zap.String("err", "This is an error"))
}

func index(resp http.ResponseWriter, req *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

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
		resp.Write(output)
		logger.Error("Index",
			zap.String("url", req.URL.Path),
			zap.String("msg", "oopen /data/company failed"),
			zap.String("err", err.Error()),
		)
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

	logger.Info("Index",
		zap.String("ip", getClientIP(req)),
		zap.Int("status", http.StatusOK),
	)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/localhost/healthz", index)
	mux.HandleFunc("/healthz", index)
	mux.HandleFunc("/error", errorHandle)
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
