package bo

import cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"

type ClusterBo struct {
	ClusterName  string
	UpstreamHost string
	UpstreamPort uint32
	ClusterType  cluster.Cluster_DiscoveryType
}
