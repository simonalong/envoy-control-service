package xds

import (
	mysqlProxyV3 "github.com/envoyproxy/go-control-plane/contrib/envoy/extensions/filters/network/mysql_proxy/v3"
	accessLogV3 "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	traceV3 "github.com/envoyproxy/go-control-plane/envoy/config/trace/v3"
	fileV3 "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	httpProxyV3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	redisProxyV3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/redis_proxy/v3"
	tcpProxyV3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	matcherV3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	_struct "github.com/golang/protobuf/ptypes/struct"
	structpb "google.golang.org/protobuf/types/known/structpb"

	"github.com/golang/protobuf/ptypes/wrappers"
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
		Address: GetListenerAddress("0.0.0.0", ListenerPort),

		// stat_prefix：用于侦听器统计信息的可选前缀
		StatPrefix: "",

		// 访问日志配置
		AccessLog: nil,

		// -------------------------------- 过滤器 --------------------------------
		// 过滤器链子
		FilterChains: FilterHttp(route),
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

func GetListenerAddress(listenerHost string, listenerPort uint32) *core.Address {
	if listenerHost == "" {
		listenerHost = "0.0.0.0"
	}
	return &core.Address{
		Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Protocol: core.SocketAddress_TCP,
				Address:  listenerHost,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: listenerPort,
				},
			},
		},
	}
}

func FilterHttp(route string) []*listener.FilterChain {
	return []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.HTTPConnectionManager,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getHttpFilter(route),
			},
		}},
	}}
}

func FilterMysql(upstreamClusterName string) []*listener.FilterChain {
	return []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.MySQLProxy,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getMysqlFilter(),
			},
		//}, {
		//	Name: wellknown.RoleBasedAccessControl,
		//	ConfigType: &listener.Filter_TypedConfig{
		//		TypedConfig: getRbacFilter(),
		//	},
		//}, {
		//	Name: wellknown.FileAccessLog,
		//	ConfigType: &listener.Filter_TypedConfig{
		//		TypedConfig: getRbacFilter(),
		//	},
		}, {
			Name: wellknown.TCPProxy,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getTcpFilter(upstreamClusterName),
			},
		}},
	}}
}

func FilterRedis(upstreamClusterName string) []*listener.FilterChain {
	return []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.RedisProxy,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getRedisFilter(upstreamClusterName),
			},
		}},
	}}
}

func FilterMongo(route, serviceName string) []*listener.FilterChain {
	return []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.MongoProxy,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getHttpFilter(route),
			},
		}},
	}}
}

func FilterTCP(route, serviceName string) []*listener.FilterChain {
	return []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.TCPProxy,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: getHttpFilter(route),
			},
		}},
	}}
}

func getMysqlFilter() *any.Any {
	mysqlProxy := &mysqlProxyV3.MySQLProxy{
		StatPrefix: "mysql",
		AccessLog:  "/var/log/envoy_egress_mysql.log",
	}

	pbst, err := anypb.New(mysqlProxy)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func getTcpFilter(upstreamCluster string) *any.Any {
	proxy := &tcpProxyV3.TcpProxy{
		StatPrefix: "tcp",
		ClusterSpecifier: &tcpProxyV3.TcpProxy_Cluster{
			Cluster: upstreamCluster,
		},
		MetadataMatch: nil,
		// 未设置则默认1小时
		IdleTimeout:           nil,
		DownstreamIdleTimeout: nil,
		UpstreamIdleTimeout:   nil,
		AccessLog: []*accessLogV3.AccessLog{{
			Name: "envoy.access_loggers.file",
			Filter: &accessLogV3.AccessLogFilter{
				FilterSpecifier: &accessLogV3.AccessLogFilter_MetadataFilter{
					MetadataFilter: &accessLogV3.MetadataFilter{
						Matcher: &matcherV3.MetadataMatcher{
							Filter: "envoy.filters.network.mysql_proxy",
							Path: []*matcherV3.MetadataMatcher_PathSegment{
								{
									Segment: &matcherV3.MetadataMatcher_PathSegment_Key{
										//Key: "control_version.biz_envoy",
										Key: "service_update.biz_envoy",
									},
								},
							},
							Value: &matcherV3.ValueMatcher{
								MatchPattern: &matcherV3.ValueMatcher_ListMatch{
									ListMatch: &matcherV3.ListMatcher{
										MatchPattern: &matcherV3.ListMatcher_OneOf{
											OneOf: &matcherV3.ValueMatcher{
												MatchPattern: &matcherV3.ValueMatcher_StringMatch{
													StringMatch: &matcherV3.StringMatcher{
														MatchPattern: &matcherV3.StringMatcher_Contains{
															// 正则表达式，记录所有操作
															Contains: "^((insert)|(INSERT)|(update)|(UPDATE)|(delete)|(DELETE)|(select)|(SELECT)|(show)|(SHOW)|(create)|(CREATE))$",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						MatchIfKeyNotFound: &wrappers.BoolValue{Value: true},
					},
				},
			},
			ConfigType: &accessLogV3.AccessLog_TypedConfig{
				TypedConfig: getAccessLogEgressMysql(),
			},
		}},

		// 最大尝试次数
		MaxConnectAttempts: &wrappers.UInt32Value{Value: 3},
		// 哈希策略
		HashPolicy: nil,
		// 用于在其他应用层上面传输tcp信息
		TunnelingConfig: nil,
	}

	pbst, err := anypb.New(proxy)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

//func getRbacFilter() *any.Any {
//	proxy := &rbacV3.RBAC{
//		StatPrefix: "rbac",
//		Rules: &rbacConfigV3.RBAC{
//			Action: rbacConfigV3.RBAC_LOG,
//			Policies: map[string]*rbacConfigV3.Policy{
//				"product-viewer": {
//					Permissions: []*rbacConfigV3.Permission{
//						{
//							Rule: &rbacConfigV3.Permission_Metadata{
//								Metadata: &matcherV3.MetadataMatcher{
//									Filter: "envoy.filters.network.mysql_proxy",
//									Path: []*matcherV3.MetadataMatcher_PathSegment{
//										{
//											Segment: &matcherV3.MetadataMatcher_PathSegment_Key{
//												//Key: "control_version.biz_envoy",
//												Key: "service_update.biz_envoy",
//											},
//										},
//									},
//									Value: &matcherV3.ValueMatcher{
//										MatchPattern: &matcherV3.ValueMatcher_ListMatch{
//											ListMatch: &matcherV3.ListMatcher{
//												MatchPattern: &matcherV3.ListMatcher_OneOf{
//													OneOf: &matcherV3.ValueMatcher{
//														MatchPattern: &matcherV3.ValueMatcher_StringMatch{
//															StringMatch: &matcherV3.StringMatcher{
//																MatchPattern: &matcherV3.StringMatcher_Contains{
//																	// 正则表达式，记录所有操作
//																	Contains: "^((insert)|(INSERT)|(update)|(UPDATE)|(delete)|(DELETE)|(select)|(SELECT)|(show)|(SHOW)|(create)|(CREATE))$",
//																},
//															},
//														},
//													},
//												},
//											},
//										},
//									},
//								},
//							},
//						},
//					},
//					Principals: []*rbacConfigV3.Principal{
//						{
//							Identifier: &rbacConfigV3.Principal_Any{
//								Any: true,
//							},
//						},
//					},
//				},
//			},
//		},
//		EnforcementType: rbacV3.RBAC_CONTINUOUS,
//	}
//
//	pbst, err := anypb.New(proxy)
//	if err != nil {
//		logger.Error("配置http连接失败")
//	}
//	return pbst
//}

func getHttpFilter(route string) *any.Any {
	// HTTP filter configuration
	manager := &httpProxyV3.HttpConnectionManager{
		CodecType:  httpProxyV3.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &httpProxyV3.HttpConnectionManager_Rds{
			Rds: &httpProxyV3.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*httpProxyV3.HttpFilter{{
			Name: wellknown.Router,
		}},
		GenerateRequestId: &wrappers.BoolValue{Value: true},
		Tracing: &httpProxyV3.HttpConnectionManager_Tracing{
			Provider: &traceV3.Tracing_Http{
				Name: "envoy.tracers.zipkin",
				//Name: "envoy.tracers.skywalking",
				ConfigType: &traceV3.Tracing_Http_TypedConfig{
					//TypedConfig: getgetDataCollectorOfZipkin(),
					//TypedConfig: getSkywalking(serviceName),
					TypedConfig: getDataCollectorOfCoreBack(),
				},
			},
			// CustomTags: []*tracingV3.CustomTag{},
		},
		AccessLog: []*accessLogV3.AccessLog{{
			//Name: "envoy.access_loggers.stdout",
			Name: "envoy.access_loggers.file",
			ConfigType: &accessLogV3.AccessLog_TypedConfig{
				TypedConfig: getAccessLogEgressHttp(),
			},
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func getRedisFilter(upstreamClusterName string) *any.Any {
	// HTTP filter configuration
	manager := &redisProxyV3.RedisProxy{
		StatPrefix: "http",
		Settings: &redisProxyV3.RedisProxy_ConnPoolSettings{
			OpTimeout: &duration.Duration{
				Seconds: 10,
			},
		},
		PrefixRoutes: &redisProxyV3.RedisProxy_PrefixRoutes{
			CatchAllRoute: &redisProxyV3.RedisProxy_PrefixRoutes_Route{
				Cluster: upstreamClusterName,
			},
		},
		DownstreamAuthPassword: &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineString{
				InlineString: "ZljIsysc0re123",
			},
		},
		//DownstreamAuthUsername: &corev3.DataSource{
		//	Specifier: &corev3.DataSource_InlineString{
		//		InlineString: "",
		//	},
		//},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func getAccessLogFilter() *any.Any {
	proxy := &accessLogV3.MetadataFilter{
		Matcher: &matcherV3.MetadataMatcher{
			Filter: "envoy.filters.network.mysql_proxy",
			Path: []*matcherV3.MetadataMatcher_PathSegment{
				{
					Segment: &matcherV3.MetadataMatcher_PathSegment_Key{
						//Key: "control_version.biz_envoy",
						Key: "service_update.biz_envoy",
					},
				},
			},
			Value: &matcherV3.ValueMatcher{
				MatchPattern: &matcherV3.ValueMatcher_ListMatch{
					ListMatch: &matcherV3.ListMatcher{
						MatchPattern: &matcherV3.ListMatcher_OneOf{
							OneOf: &matcherV3.ValueMatcher{
								MatchPattern: &matcherV3.ValueMatcher_StringMatch{
									StringMatch: &matcherV3.StringMatcher{
										MatchPattern: &matcherV3.StringMatcher_Contains{
											// 正则表达式，记录所有操作
											Contains: "^((insert)|(INSERT)|(update)|(UPDATE)|(delete)|(DELETE)|(select)|(SELECT)|(show)|(SHOW)|(create)|(CREATE))$",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		MatchIfKeyNotFound: &wrappers.BoolValue{Value: true},
	}

	pbst, err := anypb.New(proxy)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func getTcpCluster() {

}

func getAccessLogEgressHttp() *any.Any {
	// HTTP filter configuration
	accessLog := &fileV3.FileAccessLog{
		Path: "/var/log/envoy_egress_http.log",
		AccessLogFormat: &fileV3.FileAccessLog_LogFormat{
			LogFormat: &corev3.SubstitutionFormatString{
				Format: &corev3.SubstitutionFormatString_JsonFormat{
					JsonFormat: &_struct.Struct{
						Fields: map[string]*_struct.Value{
							// "start_time":                {Kind: &structpb.Value_StringValue{StringValue: "[%START_TIME%]"}},
							"method":         {Kind: &structpb.Value_StringValue{StringValue: "%REQ(:METHOD)%"}},
							"path":           {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"}},
							"protocol":       {Kind: &structpb.Value_StringValue{StringValue: "%PROTOCOL%"}},
							"response_code":  {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_CODE%"}},
							"response_flags": {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_FLAGS%"}},
							"bytes_received": {Kind: &structpb.Value_StringValue{StringValue: "%BYTES_RECEIVED%"}},
							"bytes_sent":     {Kind: &structpb.Value_StringValue{StringValue: "%BYTES_SENT%"}},
							// "duration":                  {Kind: &structpb.Value_StringValue{StringValue: "%DURATION%"}},
							"upstream_service_time": {Kind: &structpb.Value_StringValue{StringValue: "%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"}},
							"forwarded_for":         {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-FORWARDED-FOR)%"}},
							"x_request_id":          {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-REQUEST-ID)%"}},
							"trace_id":              {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-B3-TRACEID)%"}},
							"authority":             {Kind: &structpb.Value_StringValue{StringValue: "%REQ(:AUTHORITY)%"}},
							// "upstream_host":             {Kind: &structpb.Value_StringValue{StringValue: "%UPSTREAM_HOST%"}},
							"response_code_details": {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_CODE_DETAILS%"}},
							"route_name":            {Kind: &structpb.Value_StringValue{StringValue: "%ROUTE_NAME%"}},
							// "upstream_cluster":          {Kind: &structpb.Value_StringValue{StringValue: "%UPSTREAM_CLUSTER%"}},
							// "downstream_remote_address": {Kind: &structpb.Value_StringValue{StringValue: "%DOWNSTREAM_REMOTE_ADDRESS%"}},
							"dynamic_metadata": {Kind: &structpb.Value_StringValue{StringValue: "%DYNAMIC_METADATA(envoy.filters.network.mysql_proxy)%"}},
							// "hostname":                  {Kind: &structpb.Value_StringValue{StringValue: "%HOSTNAME%"}},
						},
					},
				},
			},
		},
	}
	pbst, err := anypb.New(accessLog)
	if err != nil {
		logger.Error("配置accessLog连接失败")
	}
	return pbst
}

func getAccessLogEgressMysql() *any.Any {
	// HTTP filter configuration
	accessLog := &fileV3.FileAccessLog{
		Path: "/var/log/envoy_egress_mysql.log",
		AccessLogFormat: &fileV3.FileAccessLog_LogFormat{
			LogFormat: &corev3.SubstitutionFormatString{
				Format: &corev3.SubstitutionFormatString_JsonFormat{
					JsonFormat: &_struct.Struct{
						Fields: map[string]*_struct.Value{
							// "start_time":                {Kind: &structpb.Value_StringValue{StringValue: "[%START_TIME%]"}},
							"method":         {Kind: &structpb.Value_StringValue{StringValue: "%REQ(:METHOD)%"}},
							"path":           {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"}},
							"protocol":       {Kind: &structpb.Value_StringValue{StringValue: "%PROTOCOL%"}},
							"response_code":  {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_CODE%"}},
							"response_flags": {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_FLAGS%"}},
							"bytes_received": {Kind: &structpb.Value_StringValue{StringValue: "%BYTES_RECEIVED%"}},
							"bytes_sent":     {Kind: &structpb.Value_StringValue{StringValue: "%BYTES_SENT%"}},
							// "duration":                  {Kind: &structpb.Value_StringValue{StringValue: "%DURATION%"}},
							"upstream_service_time": {Kind: &structpb.Value_StringValue{StringValue: "%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"}},
							"forwarded_for":         {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-FORWARDED-FOR)%"}},
							"x_request_id":          {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-REQUEST-ID)%"}},
							"trace_id":              {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-B3-TRACEID)%"}},
							"authority":             {Kind: &structpb.Value_StringValue{StringValue: "%REQ(:AUTHORITY)%"}},
							// "upstream_host":             {Kind: &structpb.Value_StringValue{StringValue: "%UPSTREAM_HOST%"}},
							"response_code_details": {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_CODE_DETAILS%"}},
							"route_name":            {Kind: &structpb.Value_StringValue{StringValue: "%ROUTE_NAME%"}},
							// "upstream_cluster":          {Kind: &structpb.Value_StringValue{StringValue: "%UPSTREAM_CLUSTER%"}},
							// "downstream_remote_address": {Kind: &structpb.Value_StringValue{StringValue: "%DOWNSTREAM_REMOTE_ADDRESS%"}},
							"dynamic_metadata": {Kind: &structpb.Value_StringValue{StringValue: "%DYNAMIC_METADATA(envoy.filters.network.mysql_proxy)%"}},
							// "hostname":                  {Kind: &structpb.Value_StringValue{StringValue: "%HOSTNAME%"}},
						},
					},
				},
			},
		},
	}
	pbst, err := anypb.New(accessLog)
	if err != nil {
		logger.Error("配置accessLog连接失败")
	}
	return pbst
}

func GetAccessLogEGress(proto string) *any.Any {
	accessLog := &fileV3.FileAccessLog{
		Path: "/var/log/envoy_egress_" + proto + ".log",
		AccessLogFormat: &fileV3.FileAccessLog_LogFormat{
			LogFormat: &corev3.SubstitutionFormatString{
				Format: &corev3.SubstitutionFormatString_JsonFormat{
					JsonFormat: &_struct.Struct{
						Fields: map[string]*_struct.Value{
							// "start_time":                {Kind: &structpb.Value_StringValue{StringValue: "[%START_TIME%]"}},
							"method":         {Kind: &structpb.Value_StringValue{StringValue: "%REQ(:METHOD)%"}},
							"path":           {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"}},
							"protocol":       {Kind: &structpb.Value_StringValue{StringValue: "%PROTOCOL%"}},
							"response_code":  {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_CODE%"}},
							"response_flags": {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_FLAGS%"}},
							"bytes_received": {Kind: &structpb.Value_StringValue{StringValue: "%BYTES_RECEIVED%"}},
							"bytes_sent":     {Kind: &structpb.Value_StringValue{StringValue: "%BYTES_SENT%"}},
							// "duration":                  {Kind: &structpb.Value_StringValue{StringValue: "%DURATION%"}},
							"upstream_service_time": {Kind: &structpb.Value_StringValue{StringValue: "%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"}},
							"forwarded_for":         {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-FORWARDED-FOR)%"}},
							"x_request_id":          {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-REQUEST-ID)%"}},
							"x_b3_trace":            {Kind: &structpb.Value_StringValue{StringValue: "%REQ(X-B3-TRACEID)%"}},
							"authority":             {Kind: &structpb.Value_StringValue{StringValue: "%REQ(:AUTHORITY)%"}},
							// "upstream_host":             {Kind: &structpb.Value_StringValue{StringValue: "%UPSTREAM_HOST%"}},
							"response_code_details": {Kind: &structpb.Value_StringValue{StringValue: "%RESPONSE_CODE_DETAILS%"}},
							"route_name":            {Kind: &structpb.Value_StringValue{StringValue: "%ROUTE_NAME%"}},
							// "upstream_cluster":          {Kind: &structpb.Value_StringValue{StringValue: "%UPSTREAM_CLUSTER%"}},
							// "downstream_remote_address": {Kind: &structpb.Value_StringValue{StringValue: "%DOWNSTREAM_REMOTE_ADDRESS%"}},
							"dynamic_metadata": {Kind: &structpb.Value_StringValue{StringValue: "%DYNAMIC_METADATA(envoy.filters.network.mysql_proxy)%"}},
							// "hostname":                  {Kind: &structpb.Value_StringValue{StringValue: "%HOSTNAME%"}},
						},
					},
				},
			},
		},
	}
	pbst, err := anypb.New(accessLog)
	if err != nil {
		logger.Error("配置accessLog连接失败")
	}
	return pbst
}

func getSkywalking(serviceName string) *any.Any {
	cfg := &traceV3.SkyWalkingConfig{
		GrpcService: &corev3.GrpcService{
			TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
					ClusterName: "cluster_sky",
				},
			},
		},
		ClientConfig: &traceV3.ClientConfig{
			ServiceName:  serviceName,
			InstanceName: serviceName,
		},
	}

	pbst, err := anypb.New(cfg)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func getgetDataCollectorOfZipkin() *any.Any {
	zip := &traceV3.ZipkinConfig{
		//这里使用zipkin进行搜集
		CollectorCluster:         "zipkin",
		CollectorEndpoint:        "/api/v2/spans",
		SharedSpanContext:        &wrappers.BoolValue{Value: false},
		CollectorEndpointVersion: traceV3.ZipkinConfig_HTTP_JSON,
	}

	pbst, err := anypb.New(zip)
	if err != nil {
		logger.Error("配置http连接失败")
	}
	return pbst
}

func getDataCollectorOfCoreBack() *any.Any {
	cfg := &traceV3.ZipkinConfig{
		//这里使用zipkin进行搜集
		CollectorCluster:         "cluster-core-back",
		CollectorEndpoint:        "/api/core/back/v1/http/spans",
		SharedSpanContext:        &wrappers.BoolValue{Value: false},
		CollectorEndpointVersion: traceV3.ZipkinConfig_HTTP_JSON,
	}

	pbst, err := anypb.New(cfg)
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
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "cluster_xds"},
				},
			}},
		},
	}
	return source
}
