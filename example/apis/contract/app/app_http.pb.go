// Code generated by https://github.com/oaago/protoc-gen-oaago. DO NOT EDIT.

package app

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
type BlogServiceHTTPServer interface {
	CreateArticle(context.Context, *Article) (*Article, error)

	GetArticles(context.Context, *GetArticlesReq) (*GetArticlesResp, error)
}
type RpcClientType struct {
	EtcdClient        *v3.Client
	EtcdAddr          string
	RemoteServiceName string
}

func RegisterBlogServiceHTTPServer(r oaa.Engine, srv BlogServiceHTTPServer) {
	s := BlogService{
		server: srv,
		router: r,
		resp:   defaultBlogServiceResp{},
	}
	s.RegisterService()
}

type BlogService struct {
	server BlogServiceHTTPServer
	router oaa.Engine
	resp   interface {
		Error(ctx *oaa.Ctx, err error)
		ParamsError(ctx *oaa.Ctx, err error)
		Success(ctx *oaa.Ctx, data interface{})
	}
}

// Resp 返回值
type defaultBlogServiceResp struct{}

func (resp defaultBlogServiceResp) response(ctx *oaa.Ctx, status, code int, msg string, data interface{}) {
	ctx.JSON(status, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Error 返回错误信息
func (resp defaultBlogServiceResp) Error(ctx *oaa.Ctx, err error) {
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
func (resp defaultBlogServiceResp) ParamsError(ctx *oaa.Ctx, err error) {
	_ = ctx.Error(err)
	resp.response(ctx, 400, 400, "参数错误", nil)
}

// Success 返回成功信息
func (resp defaultBlogServiceResp) Success(ctx *oaa.Ctx, data interface{}) {
	resp.response(ctx, 200, 0, "成功", data)
}

func (s *BlogService) GetArticles_0(ctx *oaa.Ctx) {
	var in GetArticlesReq

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
	out, err := s.server.(BlogServiceHTTPServer).GetArticles(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}

func (s *BlogService) GetArticles_1(ctx *oaa.Ctx) {
	var in GetArticlesReq

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
	out, err := s.server.(BlogServiceHTTPServer).GetArticles(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}

func (s *BlogService) CreateArticle_0(ctx *oaa.Ctx) {
	var in Article

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

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(BlogServiceHTTPServer).CreateArticle(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}

func (s *BlogService) RegisterService() {

	s.router.Handle("GET", "/v1/author/:author_id/articles", oaa.NewHandler(s.GetArticles_0))

	s.router.Handle("GET", "/v1/articles", oaa.NewHandler(s.GetArticles_1))

	s.router.Handle("POST", "/v1/author/:author_id/articles", oaa.NewHandler(s.CreateArticle_0))

}

func NewRpcBlogServiceClient(rpc RpcClientType) BlogServiceClient {
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
	BlogServiceClient := NewBlogServiceClient(conn)
	return BlogServiceClient
}
