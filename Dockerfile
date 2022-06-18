#build stage
FROM golang as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app
EXPOSE 32400

#image stage
FROM alpine:latest

WORKDIR /app

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/app/zoneinfo.zip

COPY --from=builder /app/application.yml /app/application.yml
COPY --from=builder /app/isc-envoy-control-service /app/isc-envoy-control-service
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /app

CMD ["./isc-envoy-control-service"]

