package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	baseServer "github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
)

func TestController() {
	baseServer.Post("cf", Cf)
}

// 新增 add
func testListener(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("value")

	logger.Info("收到值：key={}, value={}", key, value)
	rsp.SuccessOfStandard(c, "ok")
}

func Cf(c *gin.Context) {
	cfReq := CfReq{}
	err := isc.DataToObject(c.Request.Body, &cfReq)
	if err != nil {
		rsp.FailedOfStandard(c, 500, err.Error())
	}

	logger.Info("收到信息：{}", cfReq)
	rsp.SuccessOfStandard(c, "ok")
}

type CfReq struct {
	Key   string
	Value string
}

type EnvoyFTable1 struct {
	Group string
	Name  string
}
