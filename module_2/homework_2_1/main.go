package main

import (
	"fmt"
	"net/http"
	"os"
)

type Handler struct{}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		resp.Header()[k] = v
	}
	resp.Header().Set("x-server-version", os.Getenv("VERSION"))
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"msg":"success", "code":0}`))

	fmt.Printf("Request from %s, Resp status %d\n", req.RemoteAddr, http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/localhost/healthz", &Handler{})

	if err := http.ListenAndServe(":7878", &Handler{}); err != nil {
		panic(err)
	}
}
