package bo

type ListenerBo struct {
	//http、mysql、redis、etcd、mongo、dubbo等，默认http
	ListenerProto string

	ListenerName string
	RouteName    string
	// 对于非http的一些监听器直接配置了cds，绕过了Rds
	ClusterName    string

	ListenerHost string
	ListenerPort uint32
}
