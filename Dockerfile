FROM golang:1.17.7
WORKDIR ${GOPATH}/src/demo
COPY . ${GOPATH}/src/demo

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demo .

EXPOSE 8088
CMD ["./demo"]