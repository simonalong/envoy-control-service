package bo

type RouterBo struct {
	RouteName string
	RouteBind []RouteClusterBind
}

type RouteClusterBind struct {
	ClusterName string

	RoutePrefix string
}
