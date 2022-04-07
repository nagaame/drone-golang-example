FROM golang:1.18.0
WORKDIR ${GOPATH}/src/app
COPY . ${GOPATH}/src/app


RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

EXPOSE 8080
CMD ["./demo"]