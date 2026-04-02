.PHONY: help generate build run test clean

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

generate: ## 生成gRPC代码
	@echo "Generating gRPC code..."
	@mkdir -p gen/accounting/v1
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/accounting.proto
	@echo "Code generation completed!"

build: ## 编译项目
	@echo "Building..."
	@go build -o bin/grpc-server ./cmd/server
	@echo "Build completed!"

run: ## 运行服务
	@echo "Starting gRPC server..."
	@go run ./cmd/server/main.go

test: ## 运行测试
	@echo "Running tests..."
	@go test -v ./...

clean: ## 清理
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf gen/
	@go clean

lint: ## 代码检查
	@echo "Running linters..."
	@golangci-lint run

fmt: ## 格式化代码
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest

docker-build: ## 构建Docker镜像
	@echo "Building Docker image..."
	@docker build -t accounting-grpc-api:latest .

docker-run: ## 运行Docker容器
	@echo "Running Docker container..."
	@docker run -p 9090:9090 accounting-grpc-api:latest

all: clean generate build ## 完整构建流程
