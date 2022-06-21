package controller

import (
	clusterEnvoy "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/logger"
	baseServer "github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	"isc-envoy-control-service/dao"
	"isc-envoy-control-service/pojo/bo"
	"isc-envoy-control-service/pojo/dto"
	"isc-envoy-control-service/service"
)

func ListenerController() {
	baseServer.Put("add/a", addA)
	baseServer.Put("add/b", addB)
	baseServer.Put("add/c", addC)
	baseServer.Put("add/d", addD)
	baseServer.Put("add/e", addE)
	baseServer.Put("add/f", addF)

	baseServer.Put("add/all", addAll)
}

func addA(c *gin.Context) {
	logger.Info("给服务A添加数据层代理")
	version := dao.GetServiceVersion("biz-envoy-a")
	logger.Info("A：version=%v", version)
	service.SendEnvoyData(createA(version))
	dao.UpdateServiceVersion("biz-envoy-a", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addB(c *gin.Context) {
	logger.Info("给服务B添加数据层代理")
	versionB := dao.GetServiceVersion("biz-envoy-b")
	logger.Info("B：version=%v", versionB)
	service.SendEnvoyData(createB(versionB))
	dao.UpdateServiceVersion("biz-envoy-b", versionB+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addC(c *gin.Context) {
	logger.Info("给服务C添加数据层代理")
	version := dao.GetServiceVersion("biz-envoy-c")
	logger.Info("C：version=%v", version)
	service.SendEnvoyData(createC(version))
	dao.UpdateServiceVersion("biz-envoy-c", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addD(c *gin.Context) {
	logger.Info("给服务D添加数据层代理")
	version := dao.GetServiceVersion("biz-envoy-d")
	logger.Info("D：version=%v", version)
	service.SendEnvoyData(createD(version))
	dao.UpdateServiceVersion("biz-envoy-d", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addE(c *gin.Context) {
	logger.Info("给服务E添加数据层代理")
	version := dao.GetServiceVersion("biz-envoy-e")
	logger.Info("E：version=%v", version)
	service.SendEnvoyData(createE(version))
	dao.UpdateServiceVersion("biz-envoy-e", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addF(c *gin.Context) {
	logger.Info("给服务F添加数据层代理")
	versionF := dao.GetServiceVersion("biz-envoy-f")
	logger.Info("F：version=%v", versionF)
	service.SendEnvoyData(createF(versionF))
	dao.UpdateServiceVersion("biz-envoy-f", versionF+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addAll(c *gin.Context) {
	logger.Info("给所有的服务添加")

	addA(c)
	addB(c)
	addC(c)
	addD(c)
	addE(c)
	addF(c)

	rsp.SuccessOfStandard(c, "ok")
}

func createA(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-envoy-a",
		Version:     version,

		Listeners: []bo.ListenerBo{{
			ListenerName: "listener_egress_d",
			RouteName:    "route_d",
			ListenerPort: 18003,
		}},
		Routers: []bo.RouterBo{{
			RouteName: "route_d",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_d",
				RoutePrefix: "/api/d/",
			}},
		}},
		Clusters: []bo.ClusterBo{{
			ClusterName:  "cluster_d",
			UpstreamHost: "biz-envoy-d",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_jaeger",
			UpstreamHost: "jaeger-service",
			UpstreamPort: 9411,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createB(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-envoy-b",
		Version:     version,

		Listeners: []bo.ListenerBo{{
			ListenerName: "listener_egress_a",
			RouteName:    "route_a",
			ListenerPort: 18000,
		}, {
			ListenerName: "listener_egress_d",
			RouteName:    "route_d",
			ListenerPort: 18003,
		}, {
			ListenerName: "listener_egress_e",
			RouteName:    "route_e",
			ListenerPort: 18004,
		}, {
			ListenerName: "listener_egress_c",
			RouteName:    "route_c",
			ListenerPort: 18002,
		}},

		Routers: []bo.RouterBo{{
			RouteName: "route_a",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_a",
				RoutePrefix: "/api/a/",
			}},
		}, {
			RouteName: "route_d",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_d",
				RoutePrefix: "/api/d/",
			}},
		}, {
			RouteName: "route_e",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_e",
				RoutePrefix: "/api/e/",
			}},
		}, {
			RouteName: "route_c",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_c",
				RoutePrefix: "/api/c/",
			}},
		}},

		Clusters: []bo.ClusterBo{{
			ClusterName:  "cluster_a",
			UpstreamHost: "biz-envoy-a",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_d",
			UpstreamHost: "biz-envoy-d",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_e",
			UpstreamHost: "biz-envoy-e",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_c",
			UpstreamHost: "biz-envoy-c",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createC(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-envoy-c",
		Version:     version,

		Listeners: []bo.ListenerBo{{
			ListenerName: "listener_egress_f",
			RouteName:    "route_f",
			ListenerPort: 18005,
		}},
		Routers: []bo.RouterBo{{
			RouteName: "route_f",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_f",
				RoutePrefix: "/api/f/",
			}},
		}},
		Clusters: []bo.ClusterBo{{
			ClusterName:  "cluster_f",
			UpstreamHost: "biz-envoy-f",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createD(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-envoy-d",
		Version:     version,

		Listeners: []bo.ListenerBo{{
			ListenerName: "listener_egress_e",
			RouteName:    "route_e",
			ListenerPort: 18004,
		}},

		Routers: []bo.RouterBo{{
			RouteName: "route_e",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_e",
				RoutePrefix: "/api/e/",
			}},
		}},
		Clusters: []bo.ClusterBo{{
			ClusterName:  "cluster_e",
			UpstreamHost: "biz-envoy-e",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createE(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-envoy-e",
		Version:     version,

		Listeners: []bo.ListenerBo{{
			ListenerName: "listener_egress_f",
			RouteName:    "route_f",
			ListenerPort: 18005,
		}},
		Routers: []bo.RouterBo{{
			RouteName: "route_f",
			RouteBind: []bo.RouteClusterBind{{
				ClusterName: "cluster_f",
				RoutePrefix: "/api/f/",
			}},
		}},
		Clusters: []bo.ClusterBo{{
			ClusterName:  "cluster_f",
			UpstreamHost: "biz-envoy-f",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}

func createF(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-envoy-f",
		Version:     version,

		Listeners: []bo.ListenerBo{},
		Routers: []bo.RouterBo{},
		Clusters: []bo.ClusterBo{},
	}
}
