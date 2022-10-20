package main

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"isc-envoy-control-service/router"
	"isc-envoy-control-service/service"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"time"

	"github.com/isyscore/isc-gobase/logger"
	baseServer "github.com/isyscore/isc-gobase/server"

	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

func main() {
	ctx := context.Background()

	// 创建数据缓冲处理
	service.CacheData = cache.NewSnapshotCache(false, cache.IDHash{}, nil)

	// 启动grpc服务
	runGrpcServer(ctx, service.CacheData, config.GetValueInt("envoy.port"))

	router.Register()

	baseServer.Run()
}

func runGrpcServer(ctx context.Context, snapshotCacheData cache.SnapshotCache, port int) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	cb := &test.Callbacks{Debug: true}
	srv := server.NewServer(ctx, snapshotCacheData, cb)

	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, srv)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, srv)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, srv)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)

	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, srv)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error("启动端口失败: %d", port, err)
		return
	}

	go func() {
		logger.Info("grpc服务启动监听端口：%v", port)
		if err = grpcServer.Serve(lis); err != nil {
			logger.Error("启动grpc异常", err)
		}
	}()
}
