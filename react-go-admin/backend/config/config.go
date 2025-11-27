package config

const (
	// 服务器配置
	ServerPort = ":8080"

	// JWT 配置
	JWTSecret     = "your-secret-key-change-in-production"
	JWTExpireHour = 24 * 7 // 7 天

	// 数据库配置
	DBPath = "./data.db"
)

// GetServerPort 获取服务器端口
func GetServerPort() string {
	return ServerPort
}

// GetJWTSecret 获取 JWT 密钥
func GetJWTSecret() string {
	return JWTSecret
}

// GetJWTExpireHour 获取 JWT 过期时间（小时）
func GetJWTExpireHour() int {
	return JWTExpireHour
}

// GetDBPath 获取数据库路径
func GetDBPath() string {
	return DBPath
}
