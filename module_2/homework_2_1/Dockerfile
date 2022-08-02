FROM golang:1.17 as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /httpserver
# use go mod vendor to copy dependencies to the local directory, so no need to run go mod tidy
COPY ./httpserver .
RUN go build -o httpserver .

FROM scratch
COPY --from=builder /httpserver/httpserver /
EXPOSE 7878
ENTRYPOINT ["./httpserver"]
