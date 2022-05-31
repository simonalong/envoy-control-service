package service

import (
	"context"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/isyscore/isc-gobase/logger"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

func AddListener(cluster, id, listenerName, routeName, clusterName, upstreamHost string, upstreamPort uint32) {
	node := corev3.Node{Id: id, Cluster: cluster}
	// 创建数据缓冲处理
	snapshotCacheData := cache.NewSnapshotCache(false, cache.IDHash{}, nil)

	ctx := context.Background()
	snap, _ := cache.NewSnapshot("2",
		map[resource.Type][]types.Resource{
			resource.ListenerType: {getListener(listenerName, routeName)},
			resource.RouteType:    {getRouter(routeName, clusterName)},
			resource.ClusterType:  {getCluster(clusterName, upstreamHost, upstreamPort)},
		},
	)
	if err := snap.Consistent(); err != nil {
		logger.Error("数据持久化异常", err)
		return
	}

	if err := snapshotCacheData.SetSnapshot(ctx, node.GetId(), snap); err != nil {
		logger.Error("数据发送异常", err)
	}
}

func getListener(listenerName, routeName string) *listener.Listener {
	return &listener.Listener{
		// 监听器名称
		Name: listenerName,

		// 监听器地址，必须唯一
		Address: getListenerAddress(),

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: filter(routeName),
	}
}

func getRouter(routeName string, clusterName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		// 路由名称
		Name: routeName,

		// ----------------------------------- 虚拟主机 -----------------------------------
		// 组成路由表的虚拟主机数组
		VirtualHosts: getInnerHost(routeName, clusterName),
	}
}

func getInnerHost(routeName, clusterName string) []*route.VirtualHost {
	return []*route.VirtualHost{{
		Name:    routeName,
		Domains: []string{"*"},
		Routes: []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/api/biz/f/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}},
	}}
}

func getCluster(clusterName string, upstreamHost string, upstreamPort uint32) *cluster.Cluster {
	return &cluster.Cluster{
		// 集群名称
		Name: clusterName,
		// 控制层的连接超时时间
		ConnectTimeout: durationpb.New(5 * time.Second),

		// 集群类型，这里使用集群名字解析出来的第一个ip，算是逻辑ip
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STATIC},

		// 相当于dns，设置此选项是指定STATIC，STRICT_DNS或LOGICAL_DNS集群的成员所必需的
		LoadAssignment: makeInnerEndpoint(clusterName, upstreamHost, upstreamPort),

		// -------------------------- 负载均衡---------------------------------------
		// 选择主机时的负载均衡策略
		LbPolicy: cluster.Cluster_ROUND_ROBIN,
	}
}

func makeInnerEndpoint(clusterName, upstreamHost string, upstreamPort uint32) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			// 单独指定某些配置的负载均衡
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &corev3.Address{
							Address: &corev3.Address_SocketAddress{
								SocketAddress: &corev3.SocketAddress{
									Protocol: corev3.SocketAddress_TCP,
									Address:  upstreamHost,
									PortSpecifier: &corev3.SocketAddress_PortValue{
										PortValue: upstreamPort,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}
