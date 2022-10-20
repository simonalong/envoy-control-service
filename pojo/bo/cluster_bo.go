package bo

import cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"

type ClusterBo struct {
	// 集群对应的协议名称：http、mysql、redis、mongo、dubbo等，默认http
	ClusterProto string

	ClusterName  string
	UpstreamHost string
	UpstreamPort uint32
	ClusterType  cluster.Cluster_DiscoveryType
}
