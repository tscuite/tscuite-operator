FROM golang:1.14.1 as builder
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY global/ global/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go
FROM registry.cn-hangzhou.aliyuncs.com/tscuite/crd:static-nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
