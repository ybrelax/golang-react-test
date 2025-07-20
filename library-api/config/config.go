package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	ServerPort         string
	DBPath             string
	JWTSecret          string
	JWTExpiryHours     int
	RateLimitRPS       int
	RateLimitBurst     int
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
}

// LoadConfig 从环境变量和.env文件加载配置
func LoadConfig() (*Config, error) {
	// 加载.env文件
	_ = godotenv.Load()

	// 获取JWT过期时间，默认72小时
	jwtExpiryHours := 72
	if os.Getenv("JWT_EXPIRY_HOURS") != "" {
		val, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
		if err == nil {
			jwtExpiryHours = val
		}
	}

	// 获取限流配置，默认10 RPS，突发20
	rateLimitRPS := 10
	if os.Getenv("RATE_LIMIT_RPS") != "" {
		val, err := strconv.Atoi(os.Getenv("RATE_LIMIT_RPS"))
		if err == nil {
			rateLimitRPS = val
		}
	}

	rateLimitBurst := 20
	if os.Getenv("RATE_LIMIT_BURST") != "" {
		val, err := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST"))
		if err == nil {
			rateLimitBurst = val
		}
	}

	// 获取服务器端口，默认8080
	serverPort := "8080"
	if os.Getenv("SERVER_PORT") != "" {
		serverPort = os.Getenv("SERVER_PORT")
	}

	// 获取数据库路径，默认./library.db
	dbPath := "./library.db"
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	// 获取JWT密钥，必须在环境变量中设置
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-dev-secret-change-in-production"
	}

	return &Config{
		ServerPort:     serverPort,
		DBPath:         dbPath,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtExpiryHours,
		RateLimitRPS:   rateLimitRPS,
		RateLimitBurst: rateLimitBurst,
	}, nil
}
