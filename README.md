# Accounting gRPC API

Accounting System的gRPC接口层，提供标准化的gRPC API和客户端SDK。

## 功能特性

- ✅ **完整的gRPC接口**：账户管理、记账、查询、日切、调账等
- ✅ **多种记账模式**：同步、异步、批量
- ✅ **Protocol Buffers**：标准化的接口定义
- ✅ **gRPC反射**：支持grpcurl等工具调试
- ✅ **客户端SDK**：自动生成的Go客户端
- ✅ **依赖注入**：使用Uber Fx管理依赖

## 快速开始

### 1. 生成代码

```bash
# 安装protoc和插件
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成gRPC代码
make generate
```

### 2. 编译运行

```bash
# 编译
make build

# 运行
make run
```

服务默认监听在 `:9090`

### 3. 测试接口

使用grpcurl测试：

```bash
# 列出所有服务
grpcurl -plaintext localhost:9090 list

# 列出服务方法
grpcurl -plaintext localhost:9090 list accounting.v1.AccountingService

# 创建账户
grpcurl -plaintext -d '{
  "user_id": 100001,
  "account_type": "ACCOUNT_TYPE_USER",
  "category": "ACCOUNT_CATEGORY_ASSET",
  "currency": "CNY"
}' localhost:9090 accounting.v1.AccountingService/CreateAccount
```

## API接口

### 账户管理

- `CreateAccount` - 创建账户
- `GetAccount` - 查询账户
- `FreezeAccount` - 冻结账户
- `UnfreezeAccount` - 解冻账户

### 记账操作

- `DoubleEntryBooking` - 复式记账
- `BatchBooking` - 批量记账
- `MoneyFlow` - 资金流执行

### 查询操作

- `GetTransaction` - 查询流水
- `GetBalanceSnapshot` - 查询余额快照

### 管理操作

- `TriggerDayCut` - 触发日切
- `AdjustBalance` - 调账

## 使用示例

### Go客户端

```go
package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    accountingv1 "github.com/xiongwp/accounting-grpc-api/gen/accounting/v1"
)

func main() {
    // 连接服务器
    conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // 创建客户端
    client := accountingv1.NewAccountingServiceClient(conn)

    // 创建账户
    resp, err := client.CreateAccount(context.Background(), &accountingv1.CreateAccountRequest{
        UserId:      100001,
        AccountType: accountingv1.AccountType_ACCOUNT_TYPE_USER,
        Category:    accountingv1.AccountCategory_ACCOUNT_CATEGORY_ASSET,
        Currency:    "CNY",
    })

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("账户创建成功: %s", resp.Account.AccountNo)
}
```

### Python客户端

```python
import grpc
from gen.accounting.v1 import accounting_pb2, accounting_pb2_grpc

# 连接服务器
channel = grpc.insecure_channel('localhost:9090')
stub = accounting_pb2_grpc.AccountingServiceStub(channel)

# 创建账户
request = accounting_pb2.CreateAccountRequest(
    user_id=100001,
    account_type=accounting_pb2.ACCOUNT_TYPE_USER,
    category=accounting_pb2.ACCOUNT_CATEGORY_ASSET,
    currency='CNY'
)

response = stub.CreateAccount(request)
print(f'账户创建成功: {response.account.account_no}')
```

## 项目结构

```
accounting-grpc-api/
├── proto/                  # Proto文件定义
│   └── accounting.proto
├── gen/                    # 生成的代码
│   └── accounting/v1/
├── cmd/server/             # 服务器入口
│   └── main.go
├── internal/handler/       # gRPC Handler
│   └── accounting_handler.go
├── config/                 # 配置文件
├── Makefile               # 构建脚本
└── README.md              # 文档
```

## 依赖关系

```
accounting-grpc-api (本仓库)
    ↓ 依赖
accounting-system (核心服务)
    ↓ 提供
业务逻辑层服务
```

## 开发指南

### 添加新接口

1. 在 `proto/accounting.proto` 中添加接口定义
2. 运行 `make generate` 生成代码
3. 在 `internal/handler/` 中实现接口逻辑
4. 更新文档

### 调试技巧

```bash
# 使用grpcurl调试
grpcurl -plaintext localhost:9090 describe accounting.v1.AccountingService

# 查看请求/响应格式
grpcurl -plaintext localhost:9090 describe accounting.v1.CreateAccountRequest
```

## 配置说明

```yaml
# config/config.yaml
server:
  port: 9090

# 数据库配置（继承自accounting-system）
database:
  ...
```

## 性能优化

- 使用连接池管理gRPC连接
- 启用HTTP/2多路复用
- 设置合理的超时时间
- 使用拦截器记录日志和监控

## 监控指标

- gRPC请求总数
- 请求响应时间
- 错误率
- 并发连接数

## 相关链接

- 核心服务: [accounting-system](https://github.com/xiongwp/accounting-system)
- 管理后台: [accounting-admin-web](https://github.com/xiongwp/accounting-admin-web)
- gRPC官方文档: https://grpc.io/
- Protocol Buffers: https://protobuf.dev/

## License

MIT
