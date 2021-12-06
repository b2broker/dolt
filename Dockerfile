FROM golang:1.17.3-alpine
WORKDIR $GOPATH/src/github.com/b2broker/dolt
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -ldflags="-w -s" -o /go/bin/dolt ./cmd/dolt

FROM scratch
COPY --from=0 /go/bin/dolt /dolt
ENTRYPOINT ["/dolt"]