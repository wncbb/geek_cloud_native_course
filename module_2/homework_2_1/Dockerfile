FROM golang:1.17 as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /httpserver
COPY . .
RUN go mod init httpserver && go mod tidy &&  go build -o httpserver .

FROM scratch
COPY --from=builder /httpserver/httpserver /
EXPOSE 7878
ENTRYPOINT ["./httpserver"]
