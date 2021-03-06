FROM golang:1.17.3-alpine
WORKDIR $GOPATH/src/github.com/b2broker/dolt
ENV CGO_ENABLED=0
COPY . .
RUN go mod download
ARG TARGETOS TARGETARCH
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-arm64} \
    go build -ldflags="-w -s" -o /go/bin/dolt ./cmd/dolt
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-arm64} \
    go build -ldflags="-w -s" -o /go/bin/healthprobe ./cmd/healthprobe

FROM scratch
COPY --from=0 /go/bin/dolt /dolt
COPY --from=0 /go/bin/healthprobe /healthprobe
HEALTHCHECK CMD /healthprobe || exit 1
ENTRYPOINT ["/dolt"]