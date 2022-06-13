package xds

import endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"

func GetEndpoint() *endpoint.Endpoint {
	return &endpoint.Endpoint{

		// 上游主机地址
		Address: nil,
		// 选的健康检查配置用作健康检查器与健康检查主机联系的配置
		HealthCheckConfig: nil,
		// 此端点关联的主机名
		Hostname: "",
	}
}
