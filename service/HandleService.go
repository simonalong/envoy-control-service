package service

import (
	"context"
	accessLogV3 "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	redisProxyV3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/redis_proxy/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	any "github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/types/known/anypb"
	"isc-envoy-control-service/pojo/bo"
	"isc-envoy-control-service/pojo/dto"
	"isc-envoy-control-service/service/xds"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"google.golang.org/protobuf/types/known/durationpb"
)

var CacheData cache.SnapshotCache

func SendEnvoyData(envoyDataInsert *dto.EnvoyDataInsert) {
	node := corev3.Node{Id: envoyDataInsert.Id, Cluster: envoyDataInsert.ClusterName}

	ctx := context.Background()
	resourcesMap := map[resource.Type][]types.Resource{}

	var listeners []types.Resource
	if len(envoyDataInsert.Listeners) != 0 {
		for _, listenerBo := range envoyDataInsert.Listeners {
			li := AddListener(listenerBo, envoyDataInsert.ClusterName)
			if li != nil {
				listeners = append(listeners, li)
			}
		}
	}
	resourcesMap[resource.ListenerType] = listeners

	var routers []types.Resource
	if len(envoyDataInsert.Routers) != 0 {
		for _, info := range envoyDataInsert.Routers {
			routers = append(routers, AddRouter(info))
		}
	}
	resourcesMap[resource.RouteType] = routers

	var clusters []types.Resource
	if len(envoyDataInsert.Clusters) != 0 {
		for _, info := range envoyDataInsert.Clusters {
			clusters = append(clusters, AddCluster(info))
		}
	}
	resourcesMap[resource.ClusterType] = clusters

	var endpoints []types.Resource
	if len(envoyDataInsert.Endpoints) != 0 {
		for _, info := range envoyDataInsert.Endpoints {
			endpoints = append(endpoints, addEndpoint(info))
		}
	}
	resourcesMap[resource.EndpointType] = endpoints

	snap, _ := cache.NewSnapshot(isc.ToString(envoyDataInsert.Version), resourcesMap)
	if err := snap.Consistent(); err != nil {
		logger.Error("参数检查异常：%v", err)
		return
	}

	if err := CacheData.SetSnapshot(ctx, node.GetId(), snap); err != nil {
		logger.Error("数据发送异常：%v", err)
	}
}

func AddListener(listenerBo bo.ListenerBo, serviceName string) *listener.Listener {
	return getListener(listenerBo, serviceName)
}

func AddRouter(routeBo bo.RouterBo) *route.RouteConfiguration {
	return getRouter(routeBo)
}

func AddCluster(clusterBo bo.ClusterBo) *cluster.Cluster {
	return getCluster(clusterBo)
}

func addEndpoint(endpointBo bo.EndpointBo) *endpoint.Endpoint {
	// todo
	return nil
}

func getListener(listenerBo bo.ListenerBo, serviceName string) *listener.Listener {
	if listenerBo.ListenerProto == "" || listenerBo.ListenerProto == "http"{
		logger.Info("添加listener - http 【%v】", serviceName)
		return createHttpListener(listenerBo)
	} else if listenerBo.ListenerProto == "mysql" {
		logger.Info("添加listener - mysql 【%v】", serviceName)
		return createMysqlListener(listenerBo, listenerBo.ClusterName)
	} else if listenerBo.ListenerProto == "redis" {
		logger.Info("添加listener - redis 【%v】", serviceName)
		return createRedisListener(listenerBo, listenerBo.ClusterName)
	}
	return nil
}

func createHttpListener(listenerBo bo.ListenerBo) *listener.Listener {
	return &listener.Listener{
		// 监听器名称
		Name: listenerBo.ListenerName,

		// 监听器地址，必须唯一
		Address: xds.GetListenerAddress(listenerBo.ListenerHost, listenerBo.ListenerPort),

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: xds.FilterHttp(listenerBo.RouteName),
	}
}

func createMysqlListener(listenerBo bo.ListenerBo, upstreamClusterName string) *listener.Listener {
	return &listener.Listener{
		// 监听器名称
		Name: listenerBo.ListenerName,

		// 监听器地址，必须唯一
		Address: xds.GetListenerAddress(listenerBo.ListenerHost, listenerBo.ListenerPort),

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: xds.FilterMysql(upstreamClusterName),

		AccessLog: []*accessLogV3.AccessLog{{
			Name: "envoy.access_loggers.file",
			ConfigType: &accessLogV3.AccessLog_TypedConfig{
				TypedConfig: xds.GetAccessLogEGress("mysql"),
			},
		}},
	}
}

func createRedisListener(listenerBo bo.ListenerBo, upstreamClusterName string) *listener.Listener {
	return &listener.Listener{
		// 监听器名称
		Name: listenerBo.ListenerName,

		// 监听器地址，必须唯一
		Address: xds.GetListenerAddress(listenerBo.ListenerHost, listenerBo.ListenerPort),

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: xds.FilterRedis(upstreamClusterName),

		AccessLog: []*accessLogV3.AccessLog{{
			Name: "envoy.access_loggers.file",
			ConfigType: &accessLogV3.AccessLog_TypedConfig{
				TypedConfig: xds.GetAccessLogEGress("redis"),
			},
		}},
	}
}

func getRouter(routeBo bo.RouterBo) *route.RouteConfiguration {
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

func getCluster(clusterBo bo.ClusterBo) *cluster.Cluster {
	 cls := &cluster.Cluster{
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

	if clusterBo.ClusterProto == "redis"{
		cls.TypedExtensionProtocolOptions = map[string]*any.Any{
			wellknown.RedisProxy: getProtocolRedisExtension(),
		}
	}

	return cls
}

func getProtocolRedisExtension() *any.Any {
	mysqlProxy := &redisProxyV3.RedisProtocolOptions{
		AuthPassword: &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineString{
				InlineString: "ZljIsysc0re123",
			},
		},
	}

	pbst, err := anypb.New(mysqlProxy)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
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
