package controller

import (
	clusterEnvoy "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/gin-gonic/gin"
	baseServer "github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	"isc-envoy-control-service/pojo/bo"
	"isc-envoy-control-service/service"
)

func ListenerController() {
	baseServer.Get("listener/add/c", addListenerC)
	baseServer.Get("listener/add/f", addListenerF)
}

// 新增 add
func addListenerC(c *gin.Context) {
	clusterName := "default"
	id := "biz_c"

	ldxInfo := &bo.ListenerBo{
		ListenerName: "egress-f",
		RouteName:    "f-route",
		ListenerPort: 18005,
	}

	rdxInfo := &bo.RouterBo{
		RouteName: "f-route",
		RouteBind: []bo.RouteClusterBind{
			{ClusterName: "f-cluster", RoutePrefix: "/api/f/cf/ok"},
		},
	}

	cdxInfo := &bo.ClusterBo{
		ClusterName:  "f-cluster",
		UpstreamHost: "biz_envoy_f",
		UpstreamPort: 10000,
		ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
	}

	listenerInfo := service.AddListener(ldxInfo)
	routeInfo := service.AddRouter(rdxInfo)
	clusterInfo := service.AddCluster(cdxInfo)

	insertData := &bo.InsertBo{
		Cluster:       clusterName,
		Id:            id,
		Version:       "20",
		ListenerInfos: []types.Resource{listenerInfo},
		RouteInfos:    []types.Resource{routeInfo},
		ClusterInfos:  []types.Resource{clusterInfo},
	}

	service.Send(insertData)
	rsp.SuccessOfStandard(c, "ok")
}

func addListenerF(c *gin.Context) {
	cluster := "default"
	id := "biz_f"

	ldxInfo := &bo.ListenerBo{
		ListenerName: "ingress-f",
		RouteName:    "f-route",
		ListenerPort: 10000,
	}

	rdxInfo := &bo.RouterBo{
		RouteName: "f-route",
		RouteBind: []bo.RouteClusterBind{
			{ClusterName: "f-cluster", RoutePrefix: "/api/f/cf/ok/ok"},
		},
	}

	cdxInfo := &bo.ClusterBo{
		ClusterName:  "f-cluster",
		UpstreamHost: "127.0.0.1",
		UpstreamPort: 18005,
		ClusterType:  clusterEnvoy.Cluster_STATIC,
	}

	listenerInfo := service.AddListener(ldxInfo)
	routeInfo := service.AddRouter(rdxInfo)
	clusterInfo := service.AddCluster(cdxInfo)

	insertData := &bo.InsertBo{
		Cluster:       cluster,
		Id:            id,
		Version:       "11",
		ListenerInfos: []types.Resource{listenerInfo},
		RouteInfos:    []types.Resource{routeInfo},
		ClusterInfos:  []types.Resource{clusterInfo},
	}

	service.Send(insertData)
	rsp.SuccessOfStandard(c, "ok")
}
