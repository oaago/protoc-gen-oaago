type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}
type RpcClientType struct {
	EtcdClient *v3.Client
	EtcdAddr          string
	RemoteServiceName string
}
func Register{{ $.InterfaceName }}(r oaa.Engine, srv {{ $.InterfaceName }}) {
	s := {{.Name}}{
		server: srv,
		router:     r,
		resp: default{{$.Name}}Resp{},
	}
	s.RegisterService()
}

type {{$.Name}} struct{
	server {{ $.InterfaceName }}
	router oaa.Engine
	resp  interface {
		Error(ctx *oaa.Ctx, err error)
		ParamsError (ctx *oaa.Ctx, err error)
		Success(ctx *oaa.Ctx, data interface{})
	}
}

// Resp 返回值
type default{{$.Name}}Resp struct {}

func (resp default{{$.Name}}Resp) response(ctx *oaa.Ctx, status, code int, msg string, data interface{}) {
	ctx.JSON(status, map[string]interface{}{
		"code": code, 
		"msg": msg,
		"data": data,
	})
}

// Error 返回错误信息
func (resp default{{$.Name}}Resp) Error(ctx *oaa.Ctx, err error) {
	code := -1
	status := 500
	msg := "未知错误"
	
	if err == nil {
		msg += ", err is nil"
		resp.response(ctx, status, code, msg, nil)
		return
	}

	type iCode interface{
		HTTPCode() int
		Message() string
		Code() int
	}

	var c iCode
	if errors.As(err, &c) {
		status = c.HTTPCode()
		code = c.Code()
		msg = c.Message()
	}

	_ = ctx.Error(err)

	resp.response(ctx, status, code, msg, nil)
}

// ParamsError 参数错误
func (resp default{{$.Name}}Resp) ParamsError (ctx *oaa.Ctx, err error) {
	_ = ctx.Error(err)
	resp.response(ctx, 400, 400, "参数错误", nil)
}

// Success 返回成功信息
func (resp default{{$.Name}}Resp) Success(ctx *oaa.Ctx, data interface{}) {
	resp.response(ctx, 200, 0, "成功", data)
}

{{range .Methods}}
func (s *{{$.Name}}) {{ .HandlerName }} (ctx *oaa.Ctx) {
	var in {{.Request}}
{{if .HasPathParams }}
	if err := ctx.ShouldBindUri(&in); err != nil {
		// 获取validator.ValidationErrors类型的errors
        	errs, ok := err.(v10.ValidationErrors)
        	if !ok {
        		// 非validator.ValidationErrors类型错误直接返回
        		logx.Logger.Error(errs.Translate(translator.Trans))
        		s.resp.ParamsError(ctx, err)
        		return
        	}
		s.resp.ParamsError(ctx, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
	if err := ctx.ShouldBindQuery(&in); err != nil {
		// 获取validator.ValidationErrors类型的errors
        	errs, ok := err.(v10.ValidationErrors)
        	if !ok {
        		// 非validator.ValidationErrors类型错误直接返回
        		logx.Logger.Error(errs.Translate(translator.Trans))
        		s.resp.ParamsError(ctx, err)
        		return
        	}
		s.resp.ParamsError(ctx, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		// 获取validator.ValidationErrors类型的errors
        	errs, ok := err.(v10.ValidationErrors)
        	if !ok {
        		// 非validator.ValidationErrors类型错误直接返回
        		logx.Logger.Error(errs.Translate(translator.Trans))
        		s.resp.ParamsError(ctx, err)
        		return
        	}
		s.resp.ParamsError(ctx, err)
		return
	}
{{else}}
	if err := ctx.ShouldBind(&in); err != nil {
	// 获取validator.ValidationErrors类型的errors
    	errs, ok := err.(v10.ValidationErrors)
    	if !ok {
    		// 非validator.ValidationErrors类型错误直接返回
    		logx.Logger.Error(errs.Translate(translator.Trans))
    		s.resp.ParamsError(ctx, err)
    		return
    	}
		s.resp.ParamsError(ctx, err)
		return
	}
{{end}}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}
{{end}}

func (s *{{$.Name}}) RegisterService() {
{{range .Methods}}
s.router.Handle("{{.Method}}", "/api/{{.Path}}", oaa.NewHandler(s.{{ .HandlerName }}))
{{end}}
}

func NewRpc{{$.Name}}Client(rpc RpcClientType) {{$.Name}}Client {
	cli, err := v3.NewFromURL(rpc.EtcdAddr)
	if err != nil {
		panic("无法获取连接 etcd")
	}
	rpc.EtcdClient = cli
	builder, err := resolver.NewBuilder(rpc.EtcdClient)
	conn, err := grpc.Dial("etcd:///service/" + rpc.RemoteServiceName,
		grpc.WithResolvers(builder),
		grpc.WithInsecure())
	if err != nil {
		logx.Logger.Error(err)
	}
	{{$.Name}}Client := New{{$.Name}}Client(conn)
	return {{$.Name}}Client
}