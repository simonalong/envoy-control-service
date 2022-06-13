package xds

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

const (
	UpstreamHost = "www.envoyproxy.io"
	UpstreamPort = 80
)

func GetCluster(clusterName string) *cluster.Cluster {
	return &cluster.Cluster{
		// 集群名称
		Name: clusterName,
		// 控制层的连接超时时间
		ConnectTimeout: durationpb.New(5 * time.Second),

		// 集群类型，这里使用集群名字解析出来的第一个ip，算是逻辑ip
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},

		// -------------------------- dns解析 ---------------------------------------
		// dns的域名解析策略
		DnsLookupFamily: cluster.Cluster_AUTO,

		// dns刷新率， 如果指定了DNS刷新率且群集类型为STRICT_DNS或LOGICAL_DNS，则此值将用作群集的DNS刷新率。 如果未指定此设置，则默认值为5000毫秒。 对于除STRICT_DNS和LOGICAL_DNS之外的集群类型，将忽略此设置。
		DnsRefreshRate: nil,

		// dns失败刷新率，只有type为strict或者logic时候且指定了该失败率才会生效，否则，则使用dns刷新率
		DnsFailureRefreshRate: nil,

		// 相当于dns，设置此选项是指定STATIC，STRICT_DNS或LOGICAL_DNS集群的成员所必需的
		LoadAssignment: makeEndpoint(clusterName),

		// eds的一些配置，只有上面的cluster_type为eds的时候才有效
		EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{},

		// 用于设置集群的 DNS 刷新率的可选配置。如果该值设置为 true，则集群的 DNS 刷新率将设置为来自 DNS 解析的资源记录的 TTL
		RespectDnsTtl: false,

		// dns解析配置
		TypedDnsResolverConfig: nil,

		// -------------------------- 负载均衡---------------------------------------
		// 选择主机时的负载均衡策略
		LbPolicy: cluster.Cluster_ROUND_ROBIN,
		// LbPolicy进行选择指定的配置类型，比如上面的,Round，则该配置需要为 Cluster_RoundRobinLbConfig_
		LbConfig: nil,
		// 负载均衡子集配置
		LbSubsetConfig: nil,
		// 负载均衡配置，如果这个值配置了，则会替换LbPolicy
		LoadBalancingPolicy: nil,

		// -------------------------- 上游配置 ---------------------------------------
		// 清理上游的host的时间
		CleanupInterval: nil,

		// 上游配置
		UpstreamConfig: nil,

		// 上游绑定的配置
		UpstreamBindConfig: nil,

		// 用于上游连接的自定义传输套接字实现，可选
		TransportSocket: nil,

		// 上游游链接配置
		UpstreamConnectionOptions: nil,

		// 扩展协议，该字段用于为上游连接提供特定于扩展的协议选项。密钥应与扩展过滤器名称匹配，例如“envoy.filters.network.thrift_proxy”。有关特定选项的详细信息，请参阅扩展的文档。
		TypedExtensionProtocolOptions: nil,

		// -------------- 上游健康检查配置 --------------
		// 集群的可选活动健康检查配置。如果未指定配置，则不会进行健康检查，并且所有集群成员将始终被认为是健康的。
		HealthChecks: getHealthCheck(),
		// 上游不健康时候是否关闭所有链接
		CloseConnectionsOnHostHealthFailure: true,
		// 上游移除的时候，将不再考虑他的健康值
		IgnoreHealthOnHostRemoval: true,

		// -------------- 上游链接配置 --------------

		// 集群连接读写缓冲区大小的软限制。如果未指定，则应用实现定义的默认值 (1MiB)，单位应该是Byte
		PerConnectionBufferLimitBytes: &wrappers.UInt32Value{Value: 32768},

		// 集群的预链接配置
		PreconnectPolicy: nil,

		// 断路器，其实也是连接池的一些配置
		CircuitBreakers: nil,

		// -------------------------- 下游配置 ---------------------------------------
		// 为true，则集群将为每个下游连接使用单独的连接池
		ConnectionPoolPerDownstreamConnection: true,

		// -------------------------- 信息追踪配置 ---------------------------------------
		// 用于跟踪可选集群统计信息的配置
		TrackClusterStats: nil,

		// -------------------------- 其他配置 ---------------------------------------
		// 用于在预热时阻止集群就绪的可选配置
		WaitForWarmOnInit: nil,

		// 元数据字段：元数据字段可用于提供有关群集的其他信息
		Metadata: nil,

		// 一个（可选的）网络过滤器链，按应用过滤器的顺序列出。该链将应用于 Envoy 与该集群的上游服务器建立的所有传出连接。
		Filters: nil,
	}
}

func getHealthCheck() []*core.HealthCheck {
	return nil
}

func makeEndpoint(clusterName string) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			// 单独指定某些配置的负载均衡
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  UpstreamHost,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: UpstreamPort,
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
