package controller

import (
	"github.com/gin-gonic/gin"
	baseServer "github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
)

func TraceRouter() {
	baseServer.Get("group1/data", addTrace)
}

// 新增 add
func addTrace(c *gin.Context) {
	rsp.SuccessOfStandard(c, "ok")
}
