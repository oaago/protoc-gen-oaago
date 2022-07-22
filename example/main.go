package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oaago/cloud/logx"
	v1 "github.com/oaago/protoc-gen-oaago/example/api/app/app"
	"github.com/oaago/protoc-gen-oaago/example/api/app/app1"
	"github.com/oaago/protoc-gen-oaago/example/code"
	"github.com/oaago/server/oaa"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

type BlogServiceServer struct {
	v1.UnimplementedBlogServiceServer
}

type CccDddServer struct {
	app1.UnimplementedCccDddServer
}

func (c CccDddServer) ApiCccDddService(ctx context.Context, request *app1.CccDddRequest) (*app1.CccDddReply, error) {
	//TODO implement me
	err := request.Validate()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("success")
	return &app1.CccDddReply{Message: "11111"}, nil
}

func (s BlogServiceServer) CreateArticle(ctx context.Context, article *v1.Article) (*v1.Article, error) {
	if article.AuthorId < 1 {
		return nil, code.Errorf(http.StatusBadRequest, 400, "author id must > 0")
	}
	return article, nil
}

func (s BlogServiceServer) GetArticles(ctx context.Context, req *v1.GetArticlesReq) (*v1.GetArticlesResp, error) {
	if req.AuthorId < 0 {
		return nil, code.Errorf(http.StatusBadRequest, 400, "author id must >= 0")
	}
	return &v1.GetArticlesResp{
		Total: 1,
		Articles: []*v1.Article{
			{
				Title:    "test article: " + req.Title,
				Content:  "test",
				AuthorId: 1,
			},
		},
	}, nil
}

type alwaysPassLimiter struct{}

func (*alwaysPassLimiter) Limit() bool {
	return false
}

func main() {
	go func() {
		// 1 初始化 grpc 对象
		limiter := &alwaysPassLimiter{}
		grpcServer := grpc.NewServer(
			grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
				grpcctxtags.StreamServerInterceptor(),
				grpcopentracing.StreamServerInterceptor(),
				grpcprometheus.StreamServerInterceptor,
				grpczap.StreamServerInterceptor(logx.Logx),
				grpcauth.StreamServerInterceptor(func(ctx context.Context) (context.Context, error) {
					return nil, nil
				}),
				grpcrecovery.StreamServerInterceptor(),
				ratelimit.StreamServerInterceptor(limiter),
				grpc_validator.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
				grpcctxtags.UnaryServerInterceptor(),
				grpcopentracing.UnaryServerInterceptor(),
				grpcprometheus.UnaryServerInterceptor,
				grpczap.UnaryServerInterceptor(logx.Logx),
				grpcauth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
					return nil, nil
				}),
				grpcrecovery.UnaryServerInterceptor(),
				ratelimit.UnaryServerInterceptor(limiter),
				grpc_validator.UnaryServerInterceptor(),
			)),
		)
		// 2 注册服务
		v1.RegisterBlogServiceServer(grpcServer, new(BlogServiceServer))
		app1.RegisterCccDddServer(grpcServer, new(CccDddServer))
		// 3 创建监听
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			fmt.Println("net Listen err: ", err.Error())
			return
		}
		defer listener.Close()
		// 4 绑定服务
		fmt.Println("grpc....")
		grpcServer.Serve(listener)
	}()
	//
	e := gin.Default()
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{MarshalOptions: protojson.MarshalOptions{UseEnumNumbers: true, UseProtoNames: true}}))
	err := v1.RegisterBlogServiceGWFromEndpoint(ctx, mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})
	err = app1.RegisterCccDddGWFromEndpoint(ctx, mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}
	v1.RegisterBlogServiceHTTPServer(oaa.Engine{
		Engine: e,
	}, &BlogServiceServer{})
	app1.RegisterCccDddHTTPServer(oaa.Engine{
		Engine: e,
	}, &CccDddServer{})
	e.Run()
}
