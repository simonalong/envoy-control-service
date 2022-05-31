package router

import (
	"isc-envoy-control-service/controller"
)

type Option func()

var options []Option

func include() {
	// 链路处理
	options = append(options, controller.TraceRouter)
}

func Register() {
	include()
	for _, opt := range options {
		opt()
	}
}
