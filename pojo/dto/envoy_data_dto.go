package dto

import "isc-envoy-control-service/pojo/bo"

type EnvoyDataInsert struct {
	ClusterName string
	Id          string
	Version     uint32

	Listeners []bo.ListenerBo
	Routers   []bo.RouterBo
	Clusters  []bo.ClusterBo
	Endpoints []bo.EndpointBo
}
