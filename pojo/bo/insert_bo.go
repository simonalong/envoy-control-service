package bo

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

type InsertBo struct {
	Cluster string
	Id      string
	Version string

	//ListenerInfo *listener.Listener
	ListenerInfos []types.Resource
	//RouteInfo    *route.RouteConfiguration
	RouteInfos []types.Resource
	//ClusterInfo  *cluster.Cluster
	ClusterInfos []types.Resource
	//EndpointInfo  *endpoint.Endpoint
	EndpointInfos []types.Resource
}
