package xds

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/isyscore/isc-gobase/logger"
	"google.golang.org/protobuf/types/known/anypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	ListenerPort = 10000
)

func GetListener(listenerName string, route string) *listener.Listener {
	return &listener.Listener{
		// 监听器名称
		Name: "",

		// 监听器地址，必须唯一
		Address: GetListenerAddress(ListenerPort),

		// stat_prefix：用于侦听器统计信息的可选前缀
		StatPrefix: "",

		// 访问日志配置
		AccessLog: nil,

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: Filter(route),
		// 默认过滤器链
		DefaultFilterChain: nil,
		// 监听器过滤器
		ListenerFilters: nil,
		// 监听器过滤器超时时间
		ListenerFiltersTimeout: nil,
		// 当侦听器过滤器超时时是否应该创建连接。默认false
		ContinueOnListenerFiltersTimeout: false,

		// -------------------------------- tcp配置 --------------------------------
		// 侦听器是否应接受 TCP 快速打开 (TFO) 连接，这里是对应的长度
		TcpFastOpenQueueLength: nil,
		// tcp 侦听器的挂起连接队列可以增长到的最大长度，linux上默认为128
		TcpBacklogSize: nil,
		// 是否启动Mptcp，这个
		EnableMptcp: true,

		// -------------------------------- udp配置 --------------------------------
		// udp监听器配置
		UdpListenerConfig: nil,

		// -------------------------------- 套接字配置 --------------------------------
		// 套接字配置
		SocketOptions: nil,
		// 侦听器是否应设置为透明套接字
		Transparent: nil,
		// 侦听器是否应设置 *IP_FREEBIND* 套接字选项
		Freebind: nil,
		// 套接字重用配置，默认为true，连接重用
		EnableReusePort: nil,
		// 侦听器是否应绑定到端口。未绑定的侦听器只能接收从将 use_original_dst 设置为 true 的其他侦听器重定向的连接。默认为true
		BindToPort: &wrapperspb.BoolValue{Value: true},

		// -------------------------------- 连接配置 --------------------------------
		// 每个连接缓存限值，默认1Mb
		PerConnectionBufferLimitBytes: nil,
		// 连接平衡配置，不太懂？？？？
		ConnectionBalanceConfig: nil,
		// 如果使用 iptables 重定向连接，则代理接收它的端口可能与原始目标地址不同
		UseOriginalDst: nil,

		// -------------------------------- 其他 --------------------------------

		// 请求的拉取类型，没太懂？？后续测试下
		DrainType: listener.Listener_DEFAULT,
		// 指定流量相对于本地 Envoy 的预期方向，对于使用原始目标过滤器的侦听器，Windows 上需要此属性
		TrafficDirection: 0,
		// 用于表示一个 API 监听器，用于非代理客户端。向非代理应用程序公开的 API 类型取决于 API 侦听器的类型。设置此字段时，不应设置除名称以外的其他字段。
		ApiListener: nil,
	}
}

func GetListenerAddress(listenerPort uint32) *core.Address {
	return &core.Address{
		Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Protocol: core.SocketAddress_TCP,
				Address:  "0.0.0.0",
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: listenerPort,
				},
			},
		},
	}
}

func Filter(route string) []*listener.FilterChain {
	return []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.HTTPConnectionManager,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getPbst(route),
			},
		}},
	}}
}

func getPbst(route string) *any.Any {
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: wellknown.Router,
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}
