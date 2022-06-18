#build stage
FROM golang:1.18 as builder

RUN go version

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app
EXPOSE 32400
EXPOSE 11000

#image stage
FROM alpine:latest

WORKDIR /app

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/app/zoneinfo.zip

COPY --from=builder /app/application.yaml /app/application.yaml
COPY --from=builder /app/isc-envoy-control-service /app/server
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /app

CMD ["./server"]

