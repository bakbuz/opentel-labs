FROM golang:1.21
COPY . /usr/src/common-service/
WORKDIR /usr/src/common-service/
RUN go env -w GOPROXY=direct
RUN go install ./main.go
CMD ["/go/bin/main"]
