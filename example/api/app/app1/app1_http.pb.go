// Code generated by https://github.com/oaago/protoc-gen-oaago. DO NOT EDIT.

package app1

import (
	context "context"
	errors "errors"
	v10 "github.com/go-playground/validator/v10"
	logx "github.com/oaago/cloud/logx"
	oaa "github.com/oaago/server/oaa"
	translator "github.com/oaago/server/oaa/translator"
	v3 "go.etcd.io/etcd/client/v3"
	resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the https://github.com/oaago/protoc-gen-oaago package it is being compiled against.
// context.metadata.
//oaa.errors.
//v3.resolver.
//grpc.logx.
//translator.
//v10.
type CccDddHTTPServer interface {
	ApiCccDddService(context.Context, *CccDddRequest) (*CccDddReply, error)
}
type RpcClientType struct {
	EtcdClient        *v3.Client
	EtcdAddr          string
	RemoteServiceName string
}

func RegisterCccDddHTTPServer(r oaa.Engine, srv CccDddHTTPServer) {
	s := CccDdd{
		server: srv,
		router: r,
		resp:   defaultCccDddResp{},
	}
	s.RegisterService()
}

type CccDdd struct {
	server CccDddHTTPServer
	router oaa.Engine
	resp   interface {
		Error(ctx *oaa.Ctx, err error)
		ParamsError(ctx *oaa.Ctx, err error)
		Success(ctx *oaa.Ctx, data interface{})
	}
}

// Resp 返回值
type defaultCccDddResp struct{}

func (resp defaultCccDddResp) response(ctx *oaa.Ctx, status, code int, msg string, data interface{}) {
	ctx.JSON(status, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Error 返回错误信息
func (resp defaultCccDddResp) Error(ctx *oaa.Ctx, err error) {
	code := -1
	status := 500
	msg := "未知错误"

	if err == nil {
		msg += ", err is nil"
		resp.response(ctx, status, code, msg, nil)
		return
	}

	type iCode interface {
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
func (resp defaultCccDddResp) ParamsError(ctx *oaa.Ctx, err error) {
	_ = ctx.Error(err)
	resp.response(ctx, 400, 400, "参数错误", nil)
}

// Success 返回成功信息
func (resp defaultCccDddResp) Success(ctx *oaa.Ctx, data interface{}) {
	resp.response(ctx, 200, 0, "成功", data)
}

func (s *CccDdd) ApiCccDddService_0(ctx *oaa.Ctx) {
	var in CccDddRequest

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

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CccDddHTTPServer).ApiCccDddService(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}

func (s *CccDdd) RegisterService() {

	s.router.Handle("GET", "ccc/ddd/service", oaa.NewHandler(s.ApiCccDddService_0))

}

func NewRpcCccDddClient(rpc RpcClientType) CccDddClient {
	cli, err := v3.NewFromURL(rpc.EtcdAddr)
	if err != nil {
		panic("无法获取连接 etcd")
	}
	rpc.EtcdClient = cli
	builder, err := resolver.NewBuilder(rpc.EtcdClient)
	conn, err := grpc.Dial("etcd:///service/"+rpc.RemoteServiceName,
		grpc.WithResolvers(builder),
		grpc.WithInsecure())
	if err != nil {
		logx.Logger.Error(err)
	}
	CccDddClient := NewCccDddClient(conn)
	return CccDddClient
}