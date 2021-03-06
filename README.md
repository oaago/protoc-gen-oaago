# protoc-gen-go-gin

> 修改自 [kratos v2](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-http)

从 protobuf 文件中生成使用 gin 的 http rpc 服务

## 安装

请确保安装了以下依赖:

- [go 1.16](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)
- [powerproto](github.com/storyicon/powerproto/cmd/powerproto@latest)

注意由于使用 embed 特性，Go 版本必须大于 1.16

```bash
go install github.com/oaago/protoc-gen-oaago@latest
```

## 使用说明

例子见: [example](./example)

### proto 文件约定

默认情况下 rpc method 命名为 方法+资源，使用驼峰方式命名，生成代码时会进行映射

方法映射方式如下所示:

- `"GET", "FIND", "QUERY", "LIST", "SEARCH"`  --> GET
- `"POST", "CREATE"`  --> POST
- `"PUT", "UPDATE"`  --> PUT
- `"DELETE"`  --> DELETE

```protobuf
service BlogService {
  rpc CreateArticle(Article) returns (Article) {}
  // 生成 http 路由为 post: /article
}
```

除此之外还可以使用 google.api.http option 指定路由，可以通过添加 additional_bindings 使一个 rpc 方法对应多个路由

```protobuf
// blog service is a blog demo
service BlogService {
  rpc GetArticles(GetArticlesReq) returns (GetArticlesResp) {
    // 
    // 可以通过添加 additional_bindings 使一个 rpc 方法对应多个路由
    option (google.api.http) = {
      get: "/v1/articles"
      additional_bindings {
        get: "/v1/author/{author_id}/articles"
      }
    };
  }
}
```

### 文件生成

#### bash生成

```bash
  cd example
  protoc -I ./ -I ./contract \
  --proto_path=$GOPATH/src \
  --proto_path=${GOPATH}/pkg/mod \
  --proto_path=./contract/app \
  --govalidators_out=paths=source_relative:./apis/ \
  --go_out=paths=source_relative:./apis/ \
  --go-grpc_out=./apis --go-grpc_opt=paths=import \
  --oaago_out=./apis \
  --oaago_opt=paths=import \
  --grpc-gateway_out ./apis --grpc-gateway_opt paths=import \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt generate_unbound_methods=true \
  --grpc-gateway_opt register_func_suffix=GW \
  --grpc-gateway_opt allow_delete_body=true \
  --doc_out=./docs \
  --doc_opt=html,index.html \
  --openapiv2_out ./docs --openapiv2_opt logtostderr=true \
  ./contract/app/app.proto
```

#### powerproto生成

```
go install github.com/storyicon/powerproto/cmd/powerproto@latest
```

<!-- ## 相关介绍

> 待发布

- Go工程化(四) API 设计上: 项目结构 & 设计
- Go工程化(五) API 设计下: 基于 protobuf 自动生成 gin 代码 -->

