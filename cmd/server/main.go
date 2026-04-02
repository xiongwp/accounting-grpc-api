package main

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	accountingv1 "github.com/xiongwp/accounting-grpc-api/gen/accounting/v1"
	"github.com/xiongwp/accounting-grpc-api/internal/handler"
	"github.com/xiongwp/accounting-system/internal/infrastructure/database"
	"github.com/xiongwp/accounting-system/internal/infrastructure/sharding"
	"github.com/xiongwp/accounting-system/internal/repository"
	"github.com/xiongwp/accounting-system/internal/service"
)

func main() {
	app := fx.New(
		// 提供配置和日志
		fx.Provide(
			NewConfig,
			NewLogger,
		),
		// 提供基础设施
		fx.Provide(
			NewDatabaseManager,
			NewShardingRouter,
		),
		// 提供仓储层
		fx.Provide(
			repository.NewAccountRepository,
			repository.NewTransactionRepository,
		),
		// 提供服务层
		fx.Provide(
			service.NewAccountingService,
			service.NewAccountingFacadeService,
			service.NewMoneyFlowEngine,
			service.NewDayCutService,
			service.NewAdjustmentService,
		),
		// 提供gRPC Handler
		fx.Provide(
			handler.NewAccountingHandler,
		),
		// 提供gRPC Server
		fx.Provide(
			NewGRPCServer,
		),
		// 启动服务
		fx.Invoke(StartServer),
	)

	app.Run()
}

// NewConfig 创建配置
func NewConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

// NewLogger 创建日志
func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(v *viper.Viper, logger *zap.Logger) (*database.Manager, error) {
	var dbConfigs []database.DBConfig
	if err := v.UnmarshalKey("database.databases", &dbConfigs); err != nil {
		logger.Error("unmarshal database config failed", zap.Error(err))
		return nil, err
	}

	manager, err := database.NewManager(dbConfigs)
	if err != nil {
		logger.Error("create database manager failed", zap.Error(err))
		return nil, err
	}

	logger.Info("database manager created", zap.Int("dbCount", manager.DBCount()))
	return manager, nil
}

// NewShardingRouter 创建分片路由
func NewShardingRouter() *sharding.Router {
	return sharding.NewRouter()
}

// NewGRPCServer 创建gRPC服务器
func NewGRPCServer(
	handler *handler.AccountingHandler,
	logger *zap.Logger,
) *grpc.Server {
	server := grpc.NewServer()

	// 注册服务
	accountingv1.RegisterAccountingServiceServer(server, handler)

	// 启用反射（便于调试）
	reflection.Register(server)

	logger.Info("gRPC server created")
	return server
}

// StartServer 启动服务器
func StartServer(
	lc fx.Lifecycle,
	server *grpc.Server,
	v *viper.Viper,
	logger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := v.GetInt("server.port")
			if port == 0 {
				port = 9090
			}

			addr := fmt.Sprintf(":%d", port)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			logger.Info("gRPC server starting", zap.String("addr", addr))

			go func() {
				if err := server.Serve(listener); err != nil {
					logger.Error("gRPC server failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("gRPC server stopping")
			server.GracefulStop()
			return nil
		},
	})
}
