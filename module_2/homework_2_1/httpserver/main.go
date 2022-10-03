package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"httpserver/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

var logger *zap.Logger

func images(resp http.ResponseWriter, req *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	randomInMilliSeconds := rand.Intn(2000)
	time.Sleep(time.Duration(randomInMilliSeconds) * time.Millisecond)
	resp.Write([]byte(`{"code":200}`))
}

func main() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
	metrics.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/localhost/healthz", index)
	mux.HandleFunc("/healthz", index)
	mux.HandleFunc("/error", errorHandle)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/images", images)

	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    ":7878",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ServerFailed", zap.String("err", err.Error()))
		}
		logger.Info("ServerStarted")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("ShutdownFailed", zap.String("err", err.Error()))
	}
	logger.Info("ServerExitedProperly")
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
