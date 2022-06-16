package bo

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

type InsertBo struct {
	Cluster string
	Id      string
	Version string

	ListenerInfos []types.Resource
	RouteInfos    []types.Resource
	ClusterInfos  []types.Resource
	EndpointInfos []types.Resource
}
