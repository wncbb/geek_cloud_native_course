FROM golang
ENV VERSION v0.0.1
ENV GO111MODULE on
RUN cd /go/ && mkdir httpserver && cd httpserver
COPY ./main.go /go/httpserver/
RUN cd /go/httpserver/ && go mod init httpserver && go mod tidy && go build
EXPOSE 7878/tcp

ENTRYPOINT ["/go/httpserver/httpserver"]
