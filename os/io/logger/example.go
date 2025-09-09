package logger

// 使用示例
// Example 1: 使用默认配置
/*
func main() {
	// 使用默认配置创建日志管理器
	config := logger.DefaultConfig()
	logManager, err := logger.NewManager(config)
	if err != nil {
		panic(err)
	}

	// 获取默认日志器
	log := logManager.GetLogger()

	// 使用日志器
	log.Info().Msg("Application started")
	log.Debug().Str("key", "value").Msg("Debug message")
	log.Error().Err(err).Msg("An error occurred")
	log.Infof("Formatted message: %s", "hello")
}
*/

// Example 2: 使用文件输出和自定义配置
/*
func main() {
	// 自定义配置
	config := logger.Config{
		Level:      logger.InfoLevel,
		Format:     "json",                    // JSON格式输出
		Output:     "file",                    // 输出到文件
		LogDir:     "/var/log/myapp",         // 自定义日志目录
		TimeFormat: time.RFC3339,
		NoColor:    false,
		MaxAge:     30,                        // 保留30天
		MaxSize:    100,                       // 单个文件最大100MB
		MaxBackups: 30,                        // 最多保留30个文件
		Compress:   true,                      // 压缩旧日志
	}

	logManager, err := logger.NewManager(config)
	if err != nil {
		panic(err)
	}

	// 获取不同类型的日志器
	defaultLog := logManager.GetLogger()
	serviceLog := logManager.GetServiceLogger("auth-service")
	moduleLog := logManager.GetModuleLogger("database")

	// 使用不同的日志器
	defaultLog.Info().Msg("Application initialized")
	serviceLog.Info().Msg("Auth service started")
	moduleLog.Debug().Msg("Database connection established")
}
*/

// Example 3: 在依赖注入框架中使用（如Wire）
/*
// 定义服务结构体
type UserService struct {
	logger *logger.Logger
	db     *sql.DB
}

// 创建服务时注入日志器
func NewUserService(logManager *logger.Manager, db *sql.DB) *UserService {
	return &UserService{
		logger: logManager.GetServiceLogger("user-service"),
		db:     db,
	}
}

// 在服务方法中使用日志
func (s *UserService) CreateUser(user *User) error {
	s.logger.Info().Str("username", user.Username).Msg("Creating new user")

	// 业务逻辑...

	if err != nil {
		s.logger.Error().Err(err).Str("username", user.Username).Msg("Failed to create user")
		return err
	}

	s.logger.Info().Str("username", user.Username).Msg("User created successfully")
	return nil
}

// Wire provider set
var LoggerProviderSet = wire.NewSet(
	logger.DefaultConfig,
	logger.NewManager,
)
*/

// Example 4: 使用结构化日志
/*
func main() {
	config := logger.DefaultConfig()
	config.Output = "file"
	config.Format = "json"

	logManager, _ := logger.NewManager(config)
	log := logManager.GetLogger()

	// 使用结构化字段
	log.Info().
		Str("service", "payment").
		Str("method", "ProcessPayment").
		Float64("amount", 100.50).
		Str("currency", "USD").
		Str("transaction_id", "tx-12345").
		Msg("Payment processed")

	// 使用WithFields创建带默认字段的日志器
	contextLog := log.WithFields(map[string]interface{}{
		"request_id": "req-67890",
		"user_id":    "user-123",
	})

	contextLog.Info().Msg("Request started")
	contextLog.Debug().Str("action", "validate").Msg("Validating request")
	contextLog.Info().Msg("Request completed")
}
*/

// Example 5: 在HTTP中间件中使用
/*
func LoggerMiddleware(logManager *logger.Manager) gin.HandlerFunc {
	log := logManager.GetModuleLogger("http")

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// 处理请求
		c.Next()

		// 记录请求日志
		latency := time.Since(start)
		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Msg("Request processed")
	}
}
*/
