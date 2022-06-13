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
	"isc-envoy-control-service/pojo/bo"
	"time"
)

var CacheData cache.SnapshotCache

func AddCluster(clusterBo *bo.ClusterBo) *cluster.Cluster {
	return getCluster(clusterBo)
}

func AddRouter(routeBo *bo.RouterBo) *route.RouteConfiguration {
	return getRouter(routeBo)
}

func AddListener(listenerBo *bo.ListenerBo) *listener.Listener {
	return getListener(listenerBo)
}

func Send(insertBo *bo.InsertBo) {
	node := corev3.Node{Id: insertBo.Id, Cluster: insertBo.Cluster}

	ctx := context.Background()
	resourcesMap := map[resource.Type][]types.Resource{}
	if len(insertBo.ListenerInfos) != 0 {
		resourcesMap[resource.ListenerType] = insertBo.ListenerInfos
	}

	if len(insertBo.ListenerInfos) != 0 {
		resourcesMap[resource.RouteType] = insertBo.RouteInfos
	}

	if len(insertBo.ListenerInfos) != 0 {
		resourcesMap[resource.ClusterType] = insertBo.ClusterInfos
	}

	if len(insertBo.ListenerInfos) != 0 {
		resourcesMap[resource.EndpointType] = insertBo.EndpointInfos
	}

	snap, _ := cache.NewSnapshot(insertBo.Version, resourcesMap)
	if err := snap.Consistent(); err != nil {
		logger.Error("数据持久化异常", err)
		return
	}

	if err := CacheData.SetSnapshot(ctx, node.GetId(), snap); err != nil {
		logger.Error("数据发送异常", err)
	}
}

func getListener(listenerBo *bo.ListenerBo) *listener.Listener {
	return &listener.Listener{
		// 监听器名称
		Name: listenerBo.ListenerName,

		// 监听器地址，必须唯一
		Address: getListenerAddress(listenerBo.ListenerPort),

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: filter(listenerBo.RouteName),
	}
}

func getRouter(routeBo *bo.RouterBo) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		// 路由名称
		Name: routeBo.RouteName,

		// ----------------------------------- 虚拟主机 -----------------------------------
		// 组成路由表的虚拟主机数组
		VirtualHosts: getInnerHost(routeBo.RouteName, routeBo.RouteBind),
	}
}

func getInnerHost(routeName string, routeBinds []bo.RouteClusterBind) []*route.VirtualHost {
	Routes := []*route.Route{}
	for _, bind := range routeBinds {
		Routes = append(Routes, &route.Route{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: bind.RoutePrefix,
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: bind.ClusterName,
					},
				},
			},
		})
	}

	return []*route.VirtualHost{{
		Name:    routeName,
		Domains: []string{"*"},
		Routes:  Routes,
	}}
}

func getCluster(clusterBo *bo.ClusterBo) *cluster.Cluster {
	return &cluster.Cluster{
		// 集群名称
		Name: clusterBo.ClusterName,
		// 控制层的连接超时时间
		ConnectTimeout: durationpb.New(5 * time.Second),

		// 集群类型，这里使用集群名字解析出来的第一个ip，算是逻辑ip
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},

		// 相当于dns，设置此选项是指定STATIC，STRICT_DNS或LOGICAL_DNS集群的成员所必需的
		LoadAssignment: makeInnerEndpoint(clusterBo.ClusterName, clusterBo.UpstreamHost, clusterBo.UpstreamPort),

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
