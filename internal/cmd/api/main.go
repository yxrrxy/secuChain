package main

import (
	"blockSBOM/internal/api"
	"blockSBOM/internal/api/handlers"
	"blockSBOM/internal/config"
	"blockSBOM/internal/dal/dal"
	"blockSBOM/internal/dal/dal/query"
	"blockSBOM/internal/service/auth"
	"blockSBOM/internal/service/did"
	"blockSBOM/internal/service/sbom"
	"blockSBOM/internal/service/vuln"
	"blockSBOM/pkg/utils"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func setupApp(cfg *config.Config) (*server.Hertz, error) {
	// 初始化数据库
	if err := dal.InitDB(cfg); err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %v", err)
	}

	// 初始化 Fabric 客户端
	fabricClient, err := fabric.NewFabricClient(
		cfg.Fabric.ConfigPath,
		cfg.Fabric.ChannelID,
		cfg.Fabric.ChaincodeName,
	)
	if err != nil {
		return nil, fmt.Errorf("初始化Fabric客户端失败: %v", err)
	}

	// 初始化合约客户端
	didContract := contracts.NewDIDContract(fabricClient)
	sbomContract := contracts.NewSBOMContract(fabricClient)
	vulnContract := contracts.NewVulnContract(fabricClient)

	// 初始化数据库仓库
	db := dal.GetDB()
	userRepo := query.NewUserRepository(db)
	didRepo := query.NewDIDRepository(db)
	sbomRepo := query.NewSBOMRepository(db)
	vulnRepo := query.NewVulnRepository(db)

	// 初始化 JWT
	utils.InitJWTHandler(cfg.JWT.Secret)
	jwtHandler := utils.GetJWTHandler()

	// 初始化服务
	authService := auth.NewAuthService(userRepo, jwtHandler)
	didService := did.NewDIDService(didContract, didRepo)
	sbomService := sbom.NewSBOMService(sbomContract, sbomRepo)
	vulnService := vuln.NewVulnService(vulnContract, vulnRepo)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	didHandler := handlers.NewDIDHandler(didService)
	sbomHandler := handlers.NewSBOMHandler(sbomService)
	vulnHandler := handlers.NewVulnHandler(vulnService)

	// 创建服务器
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	h := server.Default(server.WithHostPorts(address))

	// 注册路由
	api.RegisterRoutes(h, authHandler, didHandler, sbomHandler, vulnHandler)

	return h, nil
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化应用
	h, err := setupApp(cfg)
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}
	defer dal.Close()

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("正在关闭服务器...")

		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()

		if err := h.Shutdown(timeoutCtx); err != nil {
			log.Printf("服务器关闭出错: %v\n", err)
		}
		cancel()
	}()

	// 启动服务器
	go h.Spin()

	// 等待关闭信号
	<-ctx.Done()
	log.Println("服务器已关闭")
}
