package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	baseServer "github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	"isc-envoy-control-service/service"
)

func ListenerController() {
	baseServer.Get("listener/add", addListener)
}

// 新增 add
func addListener(c *gin.Context) {

	cluster := "test-cluster"
	id := "test-id"
	listenerName := "test-listener"
	routeName := "test-router"
	clusterName := "test-cluster"
	upstreamHost := "192.168.8.235"
	upstreamPort := 32400
	service.AddListener(cluster, id, listenerName, routeName, clusterName, upstreamHost, isc.ToUInt32(upstreamPort))

	rsp.SuccessOfStandard(c, "ok")
}
