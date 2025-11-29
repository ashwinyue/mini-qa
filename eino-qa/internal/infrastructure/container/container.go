
package container

import (
	"context"
	"fmt"
	"time"

	"eino-qa/internal/adapter/http"
	"eino-qa/internal/adapter/http/handler"
	"eino-qa/internal/adapter/http/middleware"
	"eino-qa/internal/domain/repository"
	"eino-qa/internal/infrastructure/ai/eino"
	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/logger"
	"eino-qa/internal/infrastructure/metrics"
	"eino-qa/internal/infrastructure/repository/memory"
	"eino-qa/internal/infrastructure/repository/milvus"
	"eino-qa/internal/infrastructure/repository/sqlite"
	"eino-qa/internal/infrastructure/tenant"
	"eino-qa/internal/usecase/chat"
	"eino-qa/internal/usecase/vector"

	"github.com/gin-gonic/gin"
	milvusClient "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/sirupsen/logrus"
)

// Container 依赖注入容器
type Container struct {
	// 配置
	Config *config.Config

	// 基础设施
	Logger           logger.Logger
	LogrusLogger     *logrus.Logger
	MetricsCollector *metrics.Collector

	// 外部服务客户端
	EinoClient   *eino.Client
	MilvusClient milvusClient.Client

	// 仓储层
	VectorRepository  repository.VectorRepository
	OrderRepository   repository.OrderRepository
	SessionRepository repository.SessionRepository

	// AI 组件
	IntentRecognizer  *eino.IntentRecognizer
	RAGRetriever      *eino.RAGRetriever
	OrderQuerier      *eino.OrderQuerier
	ResponseGenerator *eino.ResponseGenerator

	// 用例层
	ChatUseCase   chat.ChatUseCaseInterface
	VectorUseCase vector.VectorUseCaseInterface

	// HTTP 层
	ChatHandler   *handler.ChatHandler
	VectorHandler *handler.VectorHandler
	HealthHandler *handler.HealthHandler
	ModelHandler  *handler.ModelHandler

	// 中间件
	TenantMiddleware   gin.HandlerFunc
	SecurityMiddleware gin.HandlerFunc
	LoggingMiddleware  gin.HandlerFunc
	MetricsMiddleware  gin.HandlerFunc
	ErrorMiddleware    gin.HandlerFunc
	AuthMiddleware     gin.HandlerFunc

	// 多租户管理
	TenantManager       *tenant.Manager
	MilvusTenantManager *milvus.TenantManager
	DBManager           *sqlite.DBManager

	// HTTP 服务器
	Router *gin.Engine
	Server *http.Server
}

// New 创建新的依赖注入容器
func New(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config: cfg,
	}

	// 按依赖顺序初始化组件
	if err := c.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := c.initMetrics(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	if err := c.initEinoClient(); err != nil {
		return nil, fmt.Errorf("failed to initialize eino client: %w", err)
	}

	if err := c.initMilvusClient(); err != nil {
		return nil, fmt.Errorf("failed to initialize milvus client: %w", err)
	}

	if err := c.initTenantManagement(); err != nil {
		return nil, fmt.Errorf("failed to initialize tenant management: %w", err)
	}

	if err := c.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := c.initAIComponents(); err != nil {
		return nil, fmt.Errorf("failed to initialize AI components: %w", err)
	}

	if err := c.initUseCases(); err != nil {
		return nil, fmt.Errorf("failed to initialize use cases: %w", err)
	}

	if err := c.initMiddlewares(); err != nil {
		return nil, fmt.Errorf("failed to initialize middlewares: %w", err)
	}

	if err := c.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	if err := c.initHTTPServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	return c, nil
}

// initLogger 初始化日志
func (c *Container) initLogger() error {
	// 创建 logrus logger
	c.LogrusLogger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(c.Config.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	c.LogrusLogger.SetLevel(level)

	// 设置日志格式
	if c.Config.Logging.Format == "json" {
		c.LogrusLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		c.LogrusLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 创建结构化日志
	log, err := logger.New(logger.Config{
		Level:    c.Config.Logging.Level,
		Format:   c.Config.Logging.Format,
		Output:   c.Config.Logging.Output,
		FilePath: c.Config.Logging.FilePath,
	})
	if err != nil {
		return err
	}

	c.Logger = log
	c.LogrusLogger.Info("logger initialized")
	return nil
}

// initMetrics 初始化指标收集器
func (c *Container) initMetrics() error {
	c.MetricsCollector = metrics.NewCollector()
	c.LogrusLogger.Info("metrics collector initialized")
	return nil
}

// initEinoClient 初始化 Eino 客户端
func (c *Container) initEinoClient() error {
	client, err := eino.NewClient(eino.ClientConfig{
		APIKey:     c.Config.DashScope.APIKey,
		ChatModel:  c.Config.DashScope.ChatModel,
		EmbedModel: c.Config.DashScope.EmbedModel,
		MaxRetries: c.Config.DashScope.MaxRetries,
		Timeout:    c.Config.DashScope.Timeout,
	})
	if err != nil {
		return err
	}

	c.EinoClient = client
	c.LogrusLogger.Info("eino client initialized")
	return nil
}

// initMilvusClient 初始化 Milvus 客户端
func (c *Container) initMilvusClient() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := milvusClient.NewClient(ctx, milvusClient.Config{
		Address: fmt.Sprintf("%s:%d", c.Config.Milvus.Host, c.Config.Milvus.Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to milvus: %w", err)
	}

	c.MilvusClient = client
	c.LogrusLogger.Info("milvus client initialized")
	return nil
}

// initTenantManagement 初始化多租户管理
func (c *Container) initTenantManagement() error {
	// 创建 Collection 管理器
	collectionManager := milvus.NewCollectionManager(c.MilvusClient, c.LogrusLogger)

	// 创建 Milvus 租户管理器
	c.MilvusTenantManager = milvus.NewTenantManager(
		collectionManager,
		c.Config.DashScope.EmbeddingDimension,
		c.LogrusLogger,
	)

	// 创建数据库管理器
	c.DBManager = sqlite.NewDBManager(c.Config.Database.BasePath)

	// 创建统一租户管理器
	c.TenantManager = tenant.NewManager(tenant.Config{
		MilvusTenantManager: c.MilvusTenantManager,
		DBManager:           c.DBManager,
		Logger:              c.LogrusLogger,
	})

	c.LogrusLogger.Info("tenant management initialized")
	return nil
}

// initRepositories 初始化仓储层
func (c *Container) initRepositories() error {
	// 向量仓储（Milvus）
	c.VectorRepository = milvus.NewVectorRepository(
		c.MilvusClient,
		c.MilvusTenantManager,
		c.EinoClient.GetEmbedModel(),
		c.LogrusLogger,
	)

	// 订单仓储（SQLite）
	c.OrderRepository = sqlite.NewOrderRepository(c.DBManager, c.LogrusLogger)

	// 会话仓储（内存实现，可以切换到 SQLite）
	c.SessionRepository = memory.NewSessionRepository(c.Config.Session.Timeout)

	c.LogrusLogger.Info("repositories initialized")
	return nil
}

// initAIComponents 初始化 AI 组件
func (c *Container) initAIComponents() error {
	// 意图识别器
	c.IntentRecognizer = eino.NewIntentRecognizer(
		c.EinoClient.GetChatModel(),
		c.Config.Intent.ConfidenceThreshold,
		c.LogrusLogger,
	)

	// RAG 检索器
	c.RAGRetriever = eino.NewRAGRetriever(
		c.EinoClient.GetChatModel(),
		c.EinoClient.GetEmbedModel(),
		c.VectorRepository,
		c.Config.RAG.TopK,
		c.Config.RAG.ScoreThreshold,
		c.LogrusLogger,
	)

	// 订单查询器
	c.OrderQuerier = eino.NewOrderQuerier(
		c.EinoClient.GetChatModel(),
		c.OrderRepository,
		c.LogrusLogger,
	)

	// 响应生成器
	c.ResponseGenerator = eino.NewResponseGenerator(
		c.EinoClient.GetChatModel(),
		c.LogrusLogger,
	)

	c.LogrusLogger.Info("AI components initialized")
	return nil
}

// initUseCases 初始化用例层
func (c *Container) initUseCases() error {
	// 对话用例
	c.ChatUseCase = chat.NewChatUseCase(
		c.IntentRecognizer,
		c.RAGRetriever,
		c.OrderQuerier,
		c.ResponseGenerator,
		c.SessionRepository,
		c.Config.Session.Timeout,
		c.Logger,
	)

	// 向量管理用例
	c.VectorUseCase = vector.NewVectorUseCase(
		c.VectorRepository,
		c.EinoClient.GetEmbedModel(),
		c.Logger,
	)

	c.LogrusLogger.Info("use cases initialized")
	return nil
}

// initMiddlewares 初始化中间件
func (c *Container) initMiddlewares() error {
	// 租户识别中间件
	c.TenantMiddleware = middleware.TenantMiddleware()

	// 安全脱敏中间件
	securityMw := middleware.NewSecurityMiddleware(c.Config.Security.SensitiveFields)
	c.SecurityMiddleware = securityMw.Handler()

	// 日志记录中间件
	c.LoggingMiddleware = middleware.NewLoggingMiddleware(c.LogrusLogger)

	// 指标收集中间件
	c.MetricsMiddleware = middleware.NewMetricsMiddleware(c.MetricsCollector)

	// 错误处理中间件
	c.ErrorMiddleware = middleware.ErrorHandler()

	// API Key 认证中间件
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.Config.Security.APIKeys)

	c.LogrusLogger.Info("middlewares initialized")
	return nil
}

// initHandlers 初始化处理器
func (c *Container) initHandlers() error {
	// 对话处理器
	c.ChatHandler = handler.NewChatHandler(c.ChatUseCase)

	// 向量管理处理器
	c.VectorHandler = handler.NewVectorHandler(c.VectorUseCase)

	// 模型管理处理器
	c.ModelHandler = handler.NewModelHandler(c.EinoClient)

	// 健康检查处理器
	c.HealthHandler = handler.NewHealthHandler().
		WithMetricsProvider(c.MetricsCollector).
		WithMilvusCheck(func(ctx context.Context) error {
			// 简单的健康检查：列出 collections
			_, err := c.MilvusClient.ListCollections(ctx)
			return err
		}).
		WithDBCheck(func(ctx context.Context) error {
			// 检查默认租户的数据库连接
			db, err := c.DBManager.GetDB("default")
			if err != nil {
				return err
			}
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Ping()
		})

	c.LogrusLogger.Info("handlers initialized")
	return nil
}

// initHTTPServer 初始化 HTTP 服务器
func (c *Container) initHTTPServer() error {
	// 设置 Gin 模式
	if c.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if c.Config.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 配置路由
	routerConfig := &http.RouterConfig{
		ChatHandler:        c.ChatHandler,
		VectorHandler:      c.VectorHandler,
		HealthHandler:      c.HealthHandler,
		ModelHandler:       c.ModelHandler,
		TenantMiddleware:   c.TenantMiddleware,
		SecurityMiddleware: c.SecurityMiddleware,
		LoggingMiddleware:  c.LoggingMiddleware,
		MetricsMiddleware:  c.MetricsMiddleware,
		ErrorMiddleware:    c.ErrorMiddleware,
		AuthMiddleware:     c.AuthMiddleware,
		Mode:               c.Config.Server.Mode,
	}

	c.Router = http.SetupRouter(routerConfig)

	// 创建服务器
	serverConfig := &http.ServerConfig{
		Host:            "0.0.0.0",
		Port:            c.Config.Server.Port,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1 MB
	}

	c.Server = http.NewServer(c.Router, serverConfig)

	c.LogrusLogger.WithField("port", c.Config.Server.Port).Info("HTTP server initialized")
	return nil
}

// Close 关闭容器，释放所有资源
func (c *Container) Close() error {
	c.LogrusLogger.Info("closing container...")

	var errs []error

	// 关闭租户管理器
	if c.TenantManager != nil {
		if err := c.TenantManager.Close(); err != nil {
			c.LogrusLogger.WithError(err).Error("failed to close tenant manager")
			errs = append(errs, err)
		}
	}

	// 关闭 Milvus 客户端
	if c.MilvusClient != nil {
		if err := c.MilvusClient.Close(); err != nil {
			c.LogrusLogger.WithError(err).Error("failed to close milvus client")
			errs = append(errs, err)
		}
	}

	// 关闭 Eino 客户端
	if c.EinoClient != nil {
		if err := c.EinoClient.Close(); err != nil {
			c.LogrusLogger.WithError(err).Error("failed to close eino client")
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing container: %v", errs)
	}

	c.LogrusLogger.Info("container closed successfully")
	return nil
}