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
	version := dao.GetServiceVersion("biz-a-service")
	logger.Info("A：version=%v", version)
	service.SendEnvoyData(createA(version))
	dao.UpdateServiceVersion("biz-a-service", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addB(c *gin.Context) {
	logger.Info("给服务B添加数据层代理")
	versionB := dao.GetServiceVersion("biz-b-service")
	logger.Info("B：version=%v", versionB)
	service.SendEnvoyData(createB(versionB))
	dao.UpdateServiceVersion("biz-b-service", versionB+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addC(c *gin.Context) {
	logger.Info("给服务C添加数据层代理")
	version := dao.GetServiceVersion("biz-c-service")
	logger.Info("C：version=%v", version)
	service.SendEnvoyData(createC(version))
	dao.UpdateServiceVersion("biz-c-service", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addD(c *gin.Context) {
	logger.Info("给服务D添加数据层代理")
	version := dao.GetServiceVersion("biz-d-service")
	logger.Info("D：version=%v", version)
	service.SendEnvoyData(createD(version))
	dao.UpdateServiceVersion("biz-d-service", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addE(c *gin.Context) {
	logger.Info("给服务E添加数据层代理")
	version := dao.GetServiceVersion("biz-e-service")
	logger.Info("E：version=%v", version)
	service.SendEnvoyData(createE(version))
	dao.UpdateServiceVersion("biz-e-service", version+1)

	rsp.SuccessOfStandard(c, "ok")
}

func addF(c *gin.Context) {
	logger.Info("给服务F添加数据层代理")
	versionF := dao.GetServiceVersion("biz-f-service")
	logger.Info("F：version=%v", versionF)
	service.SendEnvoyData(createF(versionF))
	dao.UpdateServiceVersion("biz-f-service", versionF+1)

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
		Id:          "biz-a-service",
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
			UpstreamHost: "biz-d-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "zipkin",
			UpstreamHost: "zipkin",
			UpstreamPort: 9411,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createB(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-b-service",
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
			UpstreamHost: "biz-a-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_d",
			UpstreamHost: "biz-d-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_e",
			UpstreamHost: "biz-e-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}, {
			ClusterName:  "cluster_c",
			UpstreamHost: "biz-c-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createC(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-c-service",
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
			UpstreamHost: "biz-f-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createD(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-d-service",
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
			UpstreamHost: "biz-e-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}
func createE(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-e-service",
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
			UpstreamHost: "biz-f-service",
			UpstreamPort: 10000,
			ClusterType:  clusterEnvoy.Cluster_LOGICAL_DNS,
		}},
	}
}

func createF(version uint32) *dto.EnvoyDataInsert {
	return &dto.EnvoyDataInsert{
		ClusterName: "default",
		Id:          "biz-f-service",
		Version:     version,

		Listeners: []bo.ListenerBo{},
		Routers: []bo.RouterBo{},
		Clusters: []bo.ClusterBo{},
	}
}
