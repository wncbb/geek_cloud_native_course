FROM golang:1.17 as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /httpserver
# If use go mod vendor to copy dependencies to the local directory, then no need to run go mod tidy
COPY ./httpserver .
RUN go mod tidy && go build -o httpserver .

# FROM scratch
FROM docker.io/ubuntu:20.04
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
COPY --from=builder /httpserver/httpserver /
EXPOSE 7878
ENTRYPOINT ["/tini", "--"]
CMD ["./httpserver"]
