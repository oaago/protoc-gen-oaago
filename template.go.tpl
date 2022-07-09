// rpcClient
type RpcClientType struct {
	EtcdClient *v3.Client
	EtcdAddr          string
	RemoteServiceName string
}

func NewRpc{{$.Name}}Client(rpc RpcClientType) {{$.Name}}Client {
	cli, err := v3.NewFromURL(rpc.EtcdAddr)
	if err != nil {
		panic("无法获取连接 etcd")
	}
	rpc.EtcdClient = cli
	builder, err := resolver.NewBuilder(rpc.EtcdClient)
	conn, err := grpc.Dial("etcd:///service/" + rpc.RemoteServiceName, grpc.WithResolvers(builder), grpc.WithInsecure())
	if err != nil {
		logx.Logger.Error(err)
	}
	{{$.Name}}Client := New{{$.Name}}Client(conn)
	return {{$.Name}}Client
}