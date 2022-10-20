package xds

import (
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
)

//discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
//endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
//secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, srv)
//runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, srv)

func GetRouter(routeName string, clusterName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		// 路由名称
		Name: routeName,

		// ----------------------------------- 虚拟主机 -----------------------------------
		// 组成路由表的虚拟主机数组
		VirtualHosts: getHost(clusterName),
		// 一组虚拟主机将通过 VHDS API 动态加载
		Vhds: nil,

		// ----------------------------------- 请求和响应 -----------------------------------
		// 指定应该添加到由 HTTP 连接管理器路由的每个请求的 HTTP 标头列表。
		RequestHeadersToAdd: nil,
		// 指定应从 HTTP 连接管理器路由的每个请求中删除的 HTTP 标头列表。
		RequestHeadersToRemove: nil,
		// 指定应添加到连接管理器编码的每个响应的 HTTP 标头列表。
		ResponseHeadersToAdd: nil,
		// 指定应该从连接管理器编码的每个响应中删除的 HTTP 标头列表。
		ResponseHeadersToRemove: nil,
		// 响应的最大尺寸
		MaxDirectResponseBodySizeBytes: &wrappers.UInt32Value{Value: 4096},
		// 要允许在路由或虚拟主机级别设置覆盖为true
		MostSpecificHeaderMutationsWins: false,
		//（可选）指定连接管理器将仅视为内部的 HTTP 标头列表。
		InternalOnlyHeaders: nil,

		// ----------------------------------- 其他 -----------------------------------
		// 一个可选的布尔值，指定路由表引用的集群是否将由集群管理器验证
		ValidateClusters: &wrappers.BoolValue{Value: true},

		// 插件列表
		ClusterSpecifierPlugins: nil,
	}
}

func getHost(clusterName string) []*route.VirtualHost {
	return []*route.VirtualHost{{
		Name:    "local_service",
		Domains: []string{"*"},
		Routes: []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName,
					},
					HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
						HostRewriteLiteral: UpstreamHost,
					},
				},
			},
		}},
		RateLimits: nil,
	}}
}
