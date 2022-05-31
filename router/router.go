package router

import (
	"isc-envoy-control-service/controller"
)

type Option func()

var options []Option

func include() {
	// 链路处理
	options = append(options, controller.TraceController)
	// listener监听
	options = append(options, controller.ListenerController)
	// 测试模块
	options = append(options, controller.TestController)
}

func Register() {
	include()
	for _, opt := range options {
		opt()
	}
}
